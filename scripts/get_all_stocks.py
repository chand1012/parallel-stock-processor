import csv
import os

import arrow
from fire import Fire
from alpaca.data.historical import StockHistoricalDataClient
from alpaca.data.requests import StockBarsRequest
from alpaca.data.timeframe import TimeFrame, TimeFrameUnit
from dotenv import load_dotenv
from tqdm import tqdm

load_dotenv()

API_KEY = os.getenv("ALPACA_API_KEY")
API_SECRET = os.getenv("ALPACA_API_SECRET")

# november 1st 2012
START_TIME = arrow.get("2012-11-01T23:59:00-04:00").datetime
# november 1st 2022 at 11:59pm
END_TIME = arrow.get("2022-11-01T23:59:00-04:00").datetime

SINGLE_DAY = TimeFrame(1, TimeFrameUnit.Day)


def main(ticker_csv_path, limit=0):
    # Limit of zero means no limit
    with open(ticker_csv_path, "r") as f:
        reader = csv.reader(f)
        # filter out the header row
        # also remove if the ticker has a dollar sign
        tickers = [row[0]
                   for row in reader if row[0] != "ACT Symbol" and "$" not in row[0]]

    if not os.path.isdir("data"):
        os.mkdir("data")

    if limit > 0:
        tickers = tickers[:limit]

    client = StockHistoricalDataClient(API_KEY, API_SECRET)

    for ticker in tqdm(tickers):
        if os.path.isfile(f'data/{ticker}.csv'):
            # print(f'{ticker} already exists, skipping')
            continue
        # print(f"Getting data for {ticker}")
        try:
            req_params = StockBarsRequest(
                symbol_or_symbols=ticker, start=START_TIME, end=END_TIME, timeframe=SINGLE_DAY)
            bars = client.get_stock_bars(req_params)
            bars = bars.df
            bars.to_csv(f"data/{ticker}.csv")
        except Exception as e:
            # print(f"Error getting data for {ticker}: {e}")
            continue


if __name__ == "__main__":
    Fire(main)
