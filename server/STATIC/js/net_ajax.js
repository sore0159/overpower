(function(muleObj) {

if (!muleObj.ajax) {
    muleObj.ajax = {};
}

var ajax = muleObj.ajax;

// callbacks can contain:
//         .success .error .netError .serverError .fail
ajax.getJSEND = function(url, callbacks) {
    ajax.promiseJSEND(fetch(url, { method: 'get', credentials: 'include'}), callbacks);
};


ajax.putJSEND = function(url, obj, callbacks) {
    ajax.promiseJSEND(fetch(url, { method: 'put', credentials: 'include', body: JSON.stringify(obj) }), callbacks);
};

ajax.promiseJSEND = function(promise, callbacks) {
    if (!callbacks) {
        callbacks = {};
    }
    if (!callbacks.error) {
        callbacks.error = function(msg, data) {
            console.log(msg);
            if (data) {
                console.log("DATA:", JSON.stringify(data));
            }
        };
    }
    if (!callbacks.success) {
        callbacks.success = function(msg) {
            console.log(JSON.stringify(msg));
        };
    }
    promise.then(function(response) {
        /*
        if (!response.ok) {
            if (callbacks.netError) {
                callbacks.netError("Non-ok response recieved: "+response.statusText);
            } else if (callbacks.error) {
                callbacks.error("Non-ok response recieved: "+response.statusText);
            }
            return;
        }
       */
        response.json().then(function(jsonDat) {
            if (!jsonDat) {
                if (callbacks.netError) {
                    callbacks.netError("unparsable JSON data recieved");
                } else if (callbacks.error) {
                    callbacks.error("unparsable JSON data recieved");
                }
                return;
            }
            if (jsonDat.status === "error") {
                if (callbacks.serverError) {
                    callbacks.serverError(jsonDat.message, jsonDat.code);
                } else if (callbacks.error) {
                    callbacks.error(jsonDat.message);
                }
                return;
            }
            if (jsonDat.status === "fail") {
                if (callbacks.fail) {
                    callbacks.fail(jsonDat.data);
                } else if (callbacks.error) {
                    callbacks.error("Fail JSEND response recieved", jsonDat.data);
                }
                return;
            }
            if (jsonDat.status === "success") {
                if (callbacks.success) {
                    callbacks.success(jsonDat.data);
                }
                return;
            }
            if (callbacks.netError) {
                callbacks.netError("JSON data recieved does not match JSEND specification");
            } else if (callbacks.error) {
                callbacks.error("JSON data recieved does not match JSEND specification");
            }
            return;
        }).catch(function(error) {
            if (callbacks.netError) {
                callbacks.netError("error handling data recieved: "+error);
            } else if (callbacks.error) {
                callbacks.error("error handling data recieved: "+error);
            }
        });
    }).catch(function(error) {
        if (callbacks.netError) {
            callbacks.netError("network fetch failure: "+error);
        } else if (callbacks.error) {
            callbacks.error("network fetch failure: "+error);
        }
    });
   
};




})(muleObj);


