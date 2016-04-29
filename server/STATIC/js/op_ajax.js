(function(muleObj) {

if (!muleObj.ajax) {
    console.log("OP AJAX FAILED: REQUIRES MULEOBJ AJAX LIB");
    return;
}
if (!muleObj.overpower) {
    muleObj.overpower = {};
}

// callbacks can contain:
//         .success .error .netError .serverError .fail
//ajax.getJSEND = function(url, callbacks) {
//ajax.putJSEND = function(url, obj, callbacks) {
var ajax = muleObj.ajax;
var overpower = muleObj.overpower;

overpower.getFullView = function() {
    var url = "/overpower/json/fullviews/"+overpower.GID+"/"+overpower.FID;
    var callbacks = {};
    callbacks.success = successFV;
    ajax.getJSEND(url, callbacks);
};

function successFV(data) {
    console.log("GOT", JSON.stringify(data));
    //overpower.parse.fullview(data);
    //overpower.render.map();
}



})(muleObj);


