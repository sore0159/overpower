(function(muleObj) {

var html = muleObj.html;
if (!html) {
    console.log("OP INPUT MODULE FAILED: REQUIRES HTML MODULE");
    return;
}
if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.overpower.parse) {
    muleObj.overpower.parse = {};
}
if (!muleObj.overpower.html) {
    muleObj.overpower.html = {};
}
if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}

var parse = muleObj.overpower.parse;
var data = muleObj.overpower.data;
var opHTML = muleObj.overpower.html;


})(muleObj);
