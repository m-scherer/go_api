from flask import Flask, request, jsonify, make_response
from flask_cors import CORS
from sqlalchemy import create_engine
from itertools import izip
from collections import defaultdict
import numpy as np
import pandas as pd

#import custom classes
from api_custom_classes import Normalizer, ProductIDValidator, MarketIDValidator

READ_ONLY_URL = 'postgresql://read_only_user:gocode@35.165.83.56:5432/magpie'
ENGINE = create_engine(READ_ONLY_URL)

app = Flask(__name__)
CORS(app)

#High-level: Info about a given Market
@app.route('/api/v1/markets', methods=['GET'])
def markets():
    '''
    High-level: Info about a given Market
    Endpoint: api/v1/markets
    Provides:
        Lat/Long
        Market ID
        Market Name
        Products sold at market
            Product Avg Price
            Product Price Normalized to the rest of the state
    '''

    import pdb; pdb.set_trace()

    with ENGINE.connect() as conn:
        r = conn.execute('''
            SELECT id, name, lat, long
            FROM location_xref;
            ''')
        products = conn.execute('''
            SELECT s.location_xref_id AS market, p.name AS product, AVG(s.price) as mean
            FROM sales s
            JOIN product_xref p
            ON s.product_xref_id=p.id
            GROUP BY market, product
            ORDER BY market;
            ''')
        state = conn.execute('''
            SELECT p.name AS product, s.price
            FROM sales s
            JOIN product_xref p
            ON s.product_xref_id=p.id;
            ''')

    #create list of dictionaries for market locations
    cols = ['id', 'name', 'lat', 'long']
    market_locations = [dict(izip(cols, market)) for market in r.fetchall()]

    #create normalizer class objects for each product
    df = pd.DataFrame([_ for _ in state.fetchall()], columns=['product', 'price'])
    normalize_dct = dict()
    for product, data in df.groupby('product'):
        normalizer_obj = Normalizer(product, data.copy())
        normalize_dct[product] = normalizer_obj

    #create dictionary products at each market
    products_dct = defaultdict(list)
    for market_id, product, price in products.fetchall():
        products_dct[int(market_id)].append({
                'name': product,
                'avg_price': round(float(price),2),
                'color_scale': normalize_dct[product].transform(float(price))
                })

    #combine all query data to single json object.
    #add products at each market to market summary
    market_summaries = []
    for market_dct in market_locations:
        market_dct['products'] = products_dct[market_dct['id']]
        market_summaries.append(market_dct)
    return jsonify(market_summaries)

#High-level: Product Info
@app.route('/api/v1/products/<string:product_id>', methods=['GET'])
def market(product_id):
    '''
    High-level: Product Info
    Endpoint: api/v1/products/:id
    Provides:
        Highest Price in State
        Lowest Price in State
        Median Price for State
        Unit
    '''
    #validate user input
    product_id_input = ProductIDValidator(request)
    if not product_id_input.validate():
        return jsonify(success=False, errors=product_id_input.errors)

    with ENGINE.connect() as conn:
        r = conn.execute('select price, units from sales where product_xref_id = ''' + product_id + ';')

    prices_and_units = r.fetchall()
    if len(prices_and_units) == 0:
        return jsonify([])

    prices = zip(*prices_and_units)[0]
    median, maximum, minimum, unit = np.median(prices), np.max(prices), np.min(prices), prices_and_units[0][1]
    return jsonify(dict(izip(['med', 'max', 'min', 'unit'], [median, maximum, minimum, unit])))

# High-level: Market-specific pricing info about a Product
@app.route('/api/v1/products/<string:product_id>/marketprices', methods=['GET'])
def product_price_at_market(product_id):
    '''
    High-level: Market-specific pricing info about a Product
    Endpoint: api/v1/products/:id/marketprices
    Provides:
        Product ID
        Market ID
        Average Price for this product at this market
        Calculated value of relative price of this market/product (0 = state mean, 1 = state max, -1 = state min)
    '''
    #validate user input
    product_id_input = ProductIDValidator(request)
    if not product_id_input.validate():
        return jsonify(success=False, errors=product_id_input.errors)

    with ENGINE.connect() as conn:
        r = conn.execute('SELECT price, location_xref_id FROM sales WHERE product_xref_id = ' + product_id + ';')

    df = pd.DataFrame([_ for _ in r.fetchall()], columns=['price', 'market'])
    state_avg = round(df['price'].mean(),2)
    market_df = df.groupby('market').mean().round(2)
    market_df.reset_index(inplace=True)
    market_df.rename(columns={'price': 'avg_price'}, inplace=True)
    market_df['rel_price'] = market_df['avg_price'].apply(lambda x: {True: 1, False: -1}.get(x > state_avg, 0))

    cols = market_df.columns.tolist()
    return jsonify([dict(zip(cols, (int(_[0]), _[1], int(_[2])))) for _ in market_df.values])

# Detailed: Market Info
@app.route('/api/v1/markets/<string:market_id>', methods=['GET'])
def detailed_market_info(market_id):
    '''
    Detailed: Market Info
    Endpoint: api/v1/markets/:id
    Provides:
        Market Season Dates and Operating Hours
        Website
        Address
        etc.
    '''
    #validate user input
    market_id_input = MarketIDValidator(request)
    if not market_id_input.validate():
        return jsonify(success=False, errors=market_id_input.errors)

    with ENGINE.connect() as conn:
        columns_request = conn.execute('select column_name from information_schema.columns where table_name=\'location_xref\';')
        data_request = conn.execute('select * from location_xref where id = ' + market_id + ';')
    columns = [_[0] for _ in columns_request.fetchall()]
    data = data_request.fetchall()[0]
    return jsonify([dict(zip(columns, data))])

# Detailed: Market-specific pricing info about a Product
@app.route('/api/v1/markets/<string:market_id>/productprices&PRODUCT=<string:product_name>', methods=['GET'])
def market_specific_product_price(market_id, product_name):
    '''
    Detailed: Market-specific pricing info about a Product
    Endpoint: api/v1/markets/:id/productprices&PRODUCT=product_name
    Provides:
        Week of year
        Min price for market
        Max price for market
        Mean price for market
        Median price for market
        Count - number of data points
    '''
    #validate user input
    market_id_input = MarketIDValidator(request)
    if not market_id_input.validate():
        return jsonify(success=False, errors=market_id_input.errors)

    product_id = product_lookup(product_name)
    if product_id or product_id == 0:
        product_id = str(product_id)

    with ENGINE.connect() as conn:
        r = conn.execute(' \
            SELECT price, units, DATE_PART(\'month\', date), DATE_PART(\'week\', date) \
            FROM sales \
            WHERE location_xref_id = ' + market_id + \
            ' AND product_xref_id = ' + product_id + ';')
        product_query = conn.execute('\
            SELECT name \
            FROM product_xref \
            WHERE id = ' + product_id + ';')

    product_name = product_query.fetchall()[0][0].capitalize()
    #create median, min, max values by week
    cols = ['price', 'unit', 'month', 'week']

    data = r.fetchall()
    if len(data) == 0:
        return jsonify([])

    df = pd.DataFrame(data, columns = cols)
    df['week'] = df['week'].astype(int)
    df['month'] = df['month'].astype(int)

    #create stats dataframe
    df_stats = df.groupby('week')['unit'].count().to_frame()
    df_stats = df_stats.join(df.groupby('week').median(), how='left')
    df_stats = df_stats.join(df.groupby('week')['price'].mean(), how='left', rsuffix='_mean')
    df_stats = df_stats.join(df.groupby('week')['price'].min(), how='left', rsuffix='_min')
    df_stats = df_stats.join(df.groupby('week')['price'].max(), how='left', rsuffix='_max')
    df_stats.drop('month', axis=1, inplace=True)
    #join months to stats df
    df_months = df.groupby('week')['month'].agg(lambda x: x.value_counts().index[0])
    df_stats = df_stats.join(df_months, how='left')
    df_stats.reset_index(inplace=True)
    df_stats['name'] = product_name
    df_stats.rename(columns={'unit': 'count', 'price': 'median', 'price_mean': 'mean', 'price_min': 'min', 'price_max': 'max'}, inplace=True)
    return jsonify(df_stats.to_dict(orient='records'))

# product names and corresponding ids
@app.route('/api/v1/products', methods=['GET'])
def product_name_and_id():
    '''
    High-level: Product IDs
    Endpoint: api/v1/products
    Provides:
        Product Name
        Unique Product ID
    '''
    return product_lookup()

def product_lookup(product_name=None):
    '''
    INPUT: STR (optional)
    OUTPUTS: INT, None, or List

    If argument not passed then json object of product name and ids returned. If argument passed then product_id is returned or None if product_id does not exist.
    '''
    with ENGINE.connect() as conn:
        r = conn.execute('SELECT * FROM product_xref;')

    if product_name:
        product_name = product_name.lower().strip()
        product_dct = dict((_[1], int(_[0])) for _ in r.fetchall())
        return product_dct.get(product_name, None)
    else:
        cols = ['id', 'name']
        return jsonify([dict(zip(cols, [int(_[0]), _[1]])) for _ in r.fetchall()])

if __name__ == '__main__':
    app.run(
            host='0.0.0.0',
            port=5000,
            debug=False,
            threaded=True
            )
