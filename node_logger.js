

/*
let logger = require('./logger')

logger.handler= function(obj){
  let jsonMesg = JSON.stringify(obj);
  console.log('chess logger', jsonMesg);

}

logger.with('uid', 111).info('hello world');
logger.with('uid', 111).with('server', 'chess').info('hello world');
logger.with('uid', 111).with('server', '127.0.0.1').info('hello world')

{"time":"2020-05-22T15:40:50.56012434Z","msg":"/wslogger invoked","level":"info"}
{"time":"2020-05-22T15:40:51.102213144Z","msg":"error log message 10335","level":"error","data":{"uid":84}}
{"time":"2020-05-22T15:40:52.102354792Z","msg":"error log message 10336","level":"error","data":{"uid":84}}
{"time":"2020-05-22T15:40:53.102189512Z","msg":"debug log message 10337","level":"debug","data":{"uid":13564536}}
 */

function pad(n) {
    return n < 10 ? '0' + n : n
};



function ISODateString (d) {
    return d.getUTCFullYear() + '-'
        + pad(d.getUTCMonth() + 1) + '-'
        + pad(d.getUTCDate()) + 'T'
        + pad(d.getUTCHours()) + ':'
        + pad(d.getUTCMinutes()) + ':'
        + pad(d.getUTCSeconds()) + 'Z'
};

module.exports = {
    handler: function (obj) {
        let jsonMesg = JSON.stringify(obj);
        console.log('logger', jsonMesg);
    },

    with: function (name, value) {
        let handler = this.handler;
        let result = {data: {}};
        result.data[name] = value;
        result.trace = function (message) {

            result['time'] = ISODateString(new Date());
            result['level'] = 'trace';
            result['msg'] = message;
            handler(result);
        },
            result.debug = function (message) {
                result['time'] = ISODateString(new Date());
                result['level'] = 'debug';
                result['msg'] = message;
                handler(result);
            },
            result.info = function (message) {
                result['time'] = ISODateString(new Date());
                result['level'] = 'info';
                result['msg'] = message;
                handler(result);
            },
            result.warning = function (message) {
                result['time'] = ISODateString(new Date());
                result['level'] = 'warning';
                result['msg'] = message;
                handler(result);
            },
            result.error = function (message) {
                result['time'] = ISODateString(new Date());
                result['level'] = 'error';
                result['msg'] = message;
                handler(result);
            },

            result.with = function (name, value) {
                result.data[name] = value;
                return result;
            }

        return result;
    },

    trace: function (message) {
        let result = {data: {}};
        result['time'] = ISODateString(new Date());
        result['level'] = 'trace';
        result['msg'] = message;
        handler(result);
    },
    debug: function (message) {
        let result = {data: {}};
        result['time'] = ISODateString(new Date());
        result['level'] = 'debug';
        result['msg'] = message;
        handler(result);
    },
    info: function (message) {
        let result = {data: {}};
        result['time'] = ISODateString(new Date());
        result['level'] = 'info';
        result['msg'] = message;
        handler(result);
    },
    warning: function (message) {
        let result = {data: {}};
        result['time'] = ISODateString(new Date());
        result['level'] = 'warning';
        result['msg'] = message;
        handler(result);
    },
    error: function (message) {
        let result = {data: {}};
        result['time'] = ISODateString(new Date());
        result['level'] = 'error';
        result['msg'] = message;
        handler(result);
    },
};
