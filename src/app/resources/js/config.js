var config = {
    candlestick:{
        numViews: 5,
    },
    dataTable: {
        index : 0,
        value: null
    },
    sma: {
        indexes: [],
        periods: [],
        values: []
    },
    ema: {
        indexes: [],
        periods: [],
        values: []
    },
    bbands: {
        indexes: [],
        n: 20,
        k: 2,
        up: [],
        mid: [],
        down: []
    },
    ichimoku: {
        indexes: [],
        tenkan: [],
        kijun: [],
        senkouA : [],
        senkouB: [],
        chikou: []
    },
    volume: {
        index: [],
        values: []
    },
    rsi: {
        indexes: {'up': 0, 'value': 0, 'down': 0},
        period: 14,
        up: 70,
        values: [],
        down: 30
    },
    macd: {
        indexes: [],
        periods: [],
        values: []
    },
    hv: {
        indexes: [],
        periods: [],
        values: []
    },
    events: {
        indexes: [],
        values: [],
        first: null
    }
}