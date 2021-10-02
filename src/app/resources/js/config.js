var config = {
    api:{
        enable: true,
        interval: 1000 * 3
    },
    candlestick:{
        numViews: 5,
    },
    dataTable: {
        index : 0,
        value: null
    },
    sma: {
        enable: false,
        indexes: [],
        periods: [],
        values: []
    },
    ema: {
        enable: false,
        indexes: [],
        periods: [],
        values: []
    },
    bbands: {
        enable: false,
        indexes: [],
        n: 20,
        k: 2,
        up: [],
        mid: [],
        down: []
    },
    ichimoku: {
        enable: false,
        indexes: [],
        tenkan: [],
        kijun: [],
        senkouA : [],
        senkouB: [],
        chikou: []
    },
    volume: {
        enable: false,
        index: [],
        values: []
    },
    rsi: {
        enable: false,
        indexes: {'up': 0, 'value': 0, 'down': 0},
        period: 14,
        up: 70,
        values: [],
        down: 30
    },
    macd: {
        enable: false,
        indexes: [],
        periods: [],
        values: []
    },
    hv: {
        enable: false,
        indexes: [],
        periods: [],
        values: []
    },
    events: {
        enable: false,
        indexes: [],
        values: [],
        first: null
    }
}