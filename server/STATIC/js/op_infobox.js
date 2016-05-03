(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.ajax) {
    muleObj.ajax = {};
}
if (!muleObj.html) {
    muleObj.html = {};
}


if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}
if (!muleObj.overpower.infobox) {
    muleObj.overpower.infobox = {};
}

var html = muleObj.html;
var ajax = muleObj.ajax;
var overpower = muleObj.overpower;

var data = overpower.data;
var infobox = overpower.infobox;
infobox.mainTargetText = document.getElementById('maintargettext');

})(muleObj);
