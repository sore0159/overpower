function Pixel(x, y) {
    this.x = x;
    this.y = y;
}
Pixel.prototype.eq = function(p2) {
    return (this.x == p2.x && this.y == p2.y);
};
Pixel.prototype.dist = function(p2) {
    dx = pixel.x - p2.x;
    dy = pixel.y - p2.y;
    return Math.sqrt(dx*dx+dy*dy);
};
Pixel.prototype.addPixel = function(p2) {
    var p3 = new Pixel();
    p3.x = this.x + p2.x;
    p3.y = this.y + p2.y;
    return p3;
};
Pixel.prototype.goPolar = function(p2) {
    dx = r * Math.cos(theta*Math.PI);
    dy = r * Math.sin(theta*Math.PI);
};
function Hex(x, y) {
    this.x = x;
    this.y = y;
}
Hex.prototype.eq = function(p2) {
    return (this.x == p2.x && this.y == p2.y);
};
Hex.prototype.stepsTo = function(h2) {
    var h1z = -this.x - this.y;
    var h2z = -h2.x - h2.y;
    var dx = Math.abs(h2.x - this.x);
    var dy = Math.abs(h2.y - this.y);
    var dz = Math.abs(h1z - h2z);
    return (dx + dy + dz) / 2;
};
function Grid(size, x, y) {
    this.hexRad = size;
    this.offsetX = x;
    this.offsetY = y;
}

var h1 = new Hex(10, 20);
var h2 = new Hex(11, 19);
console.log(h1.stepsTo(h2));

