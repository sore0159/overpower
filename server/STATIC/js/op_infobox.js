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

var icon = document.querySelector('link[rel="shortcut icon"]');
var blink = false;
infobox.animateIcon = function() {
    var nonBlinkSource = "/static/img/yd32.ico";
    if (overpower.data.game.newTurn) {
        blink = !blink;
        if (blink) {
            icon.href = "/static/img/yd32blink.ico";
        } else {
            icon.href = nonBlinkSource;
        }
        window.setTimeout(infobox.animateIcon, 1000);
    } else {
        blink = false;
        icon.href = nonBlinkSource;
    }
};

})(muleObj);
