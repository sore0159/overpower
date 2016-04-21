/*
        ____
       /\  /\
  ____/__\/__\
 /\  /\  /\  /
/__\/__\/__\/
\  /\  / 
 \/__\/

*/
(function(muleObj) {

if (!muleObj.geometry) {
    muleObj.geometry = {};
}


function Hex(x, y) {
   if (x) {
       this.x = x;
   } else {
       this.x = 0;
   }
   if (y) {
       this.y = y;
   } else {
       this.y = 0;
   }
}
Hex.prototype.eq = function(hex) {
    return (hex.x === this.x) && (hex.y === this.y);
};
Hex.prototype.stepsTo = function(target) {
    var thisZ = - (this.x + this.y);
    var targetZ = - (target.x + target.y);
    var dx = Math.abs(this.x - target.x);
    var dy = Math.abs(this.y - target.y);
    var dz = Math.abs(thisZ - targetZ);
    return (dx + dy + dz)/2;
};

Hex.prototype.addHex = function(hex) {
    return this.add(hex.x, hex.y);
};
Hex.prototype.add = function(x, y) {
    var dx, dy;
    if (x) {
        dx = x;
    } else {
        dx = 0;
    }
    if (y) {
        dy = y;
    } else {
        dy = 0;
    }
    var pt = new Hex(this.x + dx, this.y+dy);
    return pt;
};
Hex.prototype.scale = function(s) {
    if (!s) {
        return;
    }
    var pt = new Hex(this.x * s, this.y*s);
    return pt;
};

muleObj.geometry.Hex = Hex;

var Point = muleObj.geometry.Point;
if (!Point) {
    console.log("NO POINT FUNCTION FOUND, NOT INCLUDING HEXAGON/POINT FEATURES");
    return;
}

Hex.prototype.centerPt = function() {
    var px = 1.5 * this.x;
    // 1^2 = (.5)^2 + (dy^2)
    // dy = sqrt( .75 ) =  .86602540378
    var py = (2*this.y + this.x)*0.86602540378;
    var pt = new Point(px, py);
    return pt;
};

Hex.prototype.cornerPts = function() {
    var hexTriH = 0.86602540378;
    var center = this.centerPt();
    var pts = [
        center.add(1, 0),
        center.add(0.5, hexTriH),
        center.add(-0.5, hexTriH),
        center.add(-1, 0),
        center.add(-0.5, -hexTriH),
        center.add(0.5, -hexTriH),
    ];
    return pts;
};

Point.prototype.hexAt = function() {
	//       __
	//      | \|  width =  1.5
	//box = |_/|  height = sqrt(3) = 2*hexTriH;
    //
    var hexTriH = 0.86602540378;
    var x = this.x + 0.5;
    var y = this.y + hexTriH;
    var hx = Math.floor(x / 1.5);
    //var hy = Math.floor((y-hx*hexTriH)/ (2*hexTriH));
    var hy = Math.floor(0.5*((y/hexTriH)-hx));
    var boxX = x - 1.5*hx;
    if (boxX < 1) {
        return [hx, hy];
    }
    boxX -= 1;
    var boxY = y - (2*hy+hx)*hexTriH;
    var slope = 2*hexTriH;
    if (boxY > hexTriH) {
        boxY -= 2*hexTriH;
        if (boxX*-slope < boxY) {
            return [hx+1, hy];
        }
        return [hx, hy];
    }
    if (boxX*slope > boxY) {
        return [hx+1, hy-1];
    }
    var hex = new Hex(hx, hy);
    return hex;
};

})(muleObj);
