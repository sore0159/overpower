(function(muleObj) {

if (!muleObj.geometry) {
    console.log("OP START FAILED: REQUIRES GEOMETRY");
    return;
}

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.overpower.render) {
    muleObj.overpower.render = {};
}
if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}

var data = muleObj.overpower.data;
var render = muleObj.overpower.render;

})(muleObj);
