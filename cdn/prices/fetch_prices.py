#!/usr/bin/env python3
import functools
import json
import os
import re
import sys
import time

import boto3
import requests


TRACKED_COINS = [
    'BTC', 'ETH', 'LTC', 'XRP', 'BCH',
    'BNB', 'EOS', 'XLM', 'ETC', 'TRX', 'DOGE',
    'UNI', 'DAI', 'LINK', 'DOT', 'AAVE', 'ADA', 'SHIB',
    'FTM', 'MATIC', 'AXS', 'MANA', 'SAND', 'AVAX', 'MKR',
    'GMT', 'USDC', 'SOL', 'ATOM', 'BAT', 'GRT', 'NEAR',
    'APE', 'CHZ', 'QNT', 'XMR', 'BUSD', 'ALGO',
    'HBAR', 'YFI', 'SNX', 'ENJ', 'CRV',
    'FLOW', 'WBTC', 'LDO', 'FIL', 'DYDX', 'APT', 'MASK',
    'FLR', 'LRC', 'COMP', 'BAL', 'ENS', 'SUSHI', 'LPT',
    'GLM', 'API3', 'ONE', 'DAO', 'CVC', 'NMR', 'STORJ',
    'SNT', 'SLP', 'ANT', 'ZRX', 'IMX', 'EGLD', 'BLUR',
    'T', 'CELR', 'ARB', '1INCH', 'FLOKI', 'BABYDOGE', 'NFT', 'BTTC',
    'MAGIC', 'GMX', 'TON', 'BAND', 'CVX', 'MDT', 'SSV',
    'WLD', 'OMG',
    'ILV', 'RDNT', 'JST',
    'ELF',
    'NOT',
    'WOO', 'FET', 'RSR', 'AEVO', 'PEPE', 'OM', 'BICO', 'GALA', 'GAL',
    'AGLD', 'ETHFI', 'RNDR', 'SKL', 'TRB', 'AGIX',
]


def upload_to_s3(data, bucket, filename, content_type=None, acl=None):
    if isinstance(data, str):
        data = data.encode('utf8')
        content_type = content_type or 'text/plain'
    elif isinstance(data, (dict, list)):
        data = json.dumps(data).encode('utf8')
        content_type = content_type or 'application/json'
    else:
        content_type = content_type or 'application/octet-stream'
    s3_resource = boto3.resource(
        's3',
        endpoint_url=os.environ.get('S3_ENDPOINT') or 'https://s3.ir-thr-at1.arvanstorage.ir',
        aws_access_key_id=os.environ.get('S3_ID'),
        aws_secret_access_key=os.environ.get('S3_SECRET'),
    )
    bucket = s3_resource.Bucket(bucket)
    return bucket.put_object(
        ACL=acl or 'public-read',
        Body=data,
        Key=filename,
        ContentType=content_type,
    )


@functools.lru_cache(maxsize=1)
def get_binance_closed_markets(date_called):
    assert date_called  # Used as cache expiration key
    r = requests.get('https://api.binance.com/api/v3/exchangeInfo', timeout=10)
    data = r.json()
    closed_markets = []
    for ticker in data['symbols']:
        if ticker['status'] != 'BREAK':
            continue
        symbol = ticker['symbol']
        match = re.fullmatch('(?P<src>.*)(USDT|DAI|BTC)', symbol)
        if not match or match.group('src') not in TRACKED_COINS:
            continue
        closed_markets.append(symbol)
    print('Binance Closed Markets:', closed_markets)
    return closed_markets


def fetch_binance_spot_prices():
    closed_markets = get_binance_closed_markets(date_called=time.strftime('%Y-%m-%d'))
    r = requests.get('https://api.binance.com/api/v3/ticker/price', timeout=10)
    data = r.json()
    results = []
    btc_price = None
    for ticker in data:
        symbol = ticker['symbol']
        match = re.fullmatch('(?P<src>.*)(USDT|DAI|BTC)', symbol)
        if not match or match.group('src') not in TRACKED_COINS or symbol in closed_markets:
            continue
        if symbol == 'BTTCUSDT':
            results.append({'symbol': 'BTTUSDT', 'price': ticker['price']})
        if symbol == 'BTCUSDT':
            btc_price = ticker['price']
        results.append(ticker)
    results.append({
        'symbol': 'LASTUPDATE',
        'price': str(int(time.time())),
    })
    upload_to_s3(results, 'nobitex-cdn', 'data/prices/binance-spot.json')
    print(f'BTC={btc_price}')


def fetch_binance_futures_prices():
    r = requests.get('https://fapi.binance.com/fapi/v1/ticker/price', timeout=10)
    data = r.json()
    results = []
    btc_price = None
    for ticker in data:
        symbol = ticker['symbol']
        price = ticker['price']
        if not symbol.endswith('USDT'):
            continue
        if symbol == 'BTCUSDT':
            btc_price = price
        results.append({'symbol': symbol, 'price': price})
    results.append({
        'symbol': 'LASTUPDATE',
        'price': str(int(time.time())),
    })
    upload_to_s3(results, 'nobitex-cdn', 'data/prices/binance-futures.json')
    print(f'BTC={btc_price}')


def fetch_okx_spot_prices():
    r = requests.get('https://www.okx.com/api/v5/market/tickers?instType=SPOT', timeout=10)
    data = r.json()['data']
    results = []
    btc_price = None
    for ticker in data:
        src, dst = ticker['instId'].split('-')
        if dst != 'USDT' or src not in TRACKED_COINS:
            continue
        price = ticker['last']
        if src == 'BTC':
            btc_price = price
        results.append({
            'symbol': src + dst,
            'price': price,
        })
    results.append({
        'symbol': 'LASTUPDATE',
        'price': str(int(time.time())),
    })
    upload_to_s3(results, 'nobitex-cdn', 'data/prices/okx-spot.json')
    print(f'BTC={btc_price}')


def fetch_prices(source):
    if source == 'binance-spot':
        runner = fetch_binance_spot_prices
        wait = 15
    elif source == 'binance-futures':
        runner = fetch_binance_futures_prices
        wait = 10
    elif source == 'okx-spot':
        runner = fetch_okx_spot_prices
        wait = 15
    else:
        runner = lambda: print('Unknown source!')
        wait = 60
    # Run Loop
    while True:
        try:
            runner()
        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f'Error: {e}')
        time.sleep(wait)


if __name__ == '__main__':
    fetch_prices(sys.argv[1])
