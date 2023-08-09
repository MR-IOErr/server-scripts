#!/usr/bin/env python3
import json
import time

import requests


TRACKED_COINS = ['BTC', 'ETH', 'SHIB', 'DOGE', 'BABYDOGE', 'DAO', 'DAI', 'GLM', 'QNT', 'SNT', 'NFT', 'NMR', 'CVC', 'ILV', 'SLP', 'CVX', 'ELF', 'MDT', 'BLUR', 'MAGIC', 'IMX', 'FLOKI', 'TON', 'LPT', 'API3']


def fetch_prices():
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
    with open('/var/www/rawdata1/okx/spot_price.json', 'w') as f:
        json.dump(results, f)
    print(f'BTC={btc_price}')


while True:
    try:
        fetch_prices()
    except KeyboardInterrupt:
        break
    except Exception as e:
        print(f'Error: {e}')
    time.sleep(1)
