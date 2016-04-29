(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.overpower.parse) {
    muleObj.overpower.parse = {};
}
if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}

var parse = muleObj.overpower.parse;
var data = muleObj.overpower.data;


parse.fullview = function(fullview) {
    data.fullview = fullview;
};

parse.rotateStars = function(clockwise, scale) {
    var dTheta = (clockwise) ? 1: -1;
    scale = (scale) ? scale : 1;
    dTheta *= scale * 0.01;
    data.stars.screen.rotate(dTheta);
};

parse.rotateMap = function(clockwise, scale) {
    var dTheta = (clockwise) ? 1: -1;
    scale = (scale) ? scale : 1;
    dTheta *= scale * 0.01;
    data.map.screen.rotate(dTheta);
};


})(muleObj);
