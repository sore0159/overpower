(function(muleObj) {

if (!muleObj.geometry) {
    muleObj.geometry = {};
}

function Point(x, y) {
    this.x = x || 0;
    this.y = y || 0;
}

Point.prototype.eq = function(point) {
    return (point.x === this.x) && (point.y === this.y);
};
Point.prototype.dist = function(point) {
    var dx = this.x - point.x;
    var dy = this.y - point.y;
    return Math.sqrt(dx*dx + dy*dy);
};

Point.prototype.addPoint = function(point) {
    return this.add(point.x, point.y);
};
Point.prototype.subPoint = function(point) {
    return this.add(-point.x, -point.y);
};
Point.prototype.add = function(x, y) {
    var dx = x || 0;
    var dy = y || 0;
    var pt = new Point(this.x + dx, this.y+dy);
    return pt;
};
Point.prototype.scale = function(s) {
    if (!s) {
        return;
    }
    var pt = new Point(this.x * s, this.y*s);
    return pt;
};
Point.prototype.addPolar = function(r, theta) {
    theta = Math.PI * theta;
    var dx = r * Math.cos(theta);
    var dy = r * Math.sin(theta);
    var pt = new Point(this.x+dx, this.y+dy);
    return pt;
};
Point.prototype.polarTo = function(point) {
    if (this.eq(point)) {
        return [0, 0];  // Not point values!  r, theta  -- not x,y
    }
    var r = this.dist(point);
    var cosT = (point.x - this.x) / r;
    var theta = Math.acos(cosT)/ Math.PI;
    if (point.y < this.y) {
        theta = 2 - theta;
    }
    return [r, theta];  // Not point values!  r, theta  -- not x,y
};


muleObj.geometry.Point = Point;


function Transform(scale, originX, originY, theta, squashY, noMirrorY) {
    this.scale = scale || 1;
    this.originX = originX || 0;
    this.originY = originY || 0;
    theta = theta || 0;
    this.theta = theta;
    this.mirrorY = !(noMirrorY);
    this.squashY = squashY || 1;
    this.cosTIn = Math.cos(theta*Math.PI);
    this.sinTIn = Math.sin(theta*Math.PI);
    this.cosTOut = Math.cos(-1*theta*Math.PI);
    this.sinTOut = Math.sin(-1*theta*Math.PI);
}
Transform.prototype.setTheta = function(theta) {
    this.theta = theta;
    this.cosTIn = Math.cos(theta*Math.PI);
    this.sinTIn = Math.sin(theta*Math.PI);
    this.cosTOut = Math.cos(-1*theta*Math.PI);
    this.sinTOut = Math.sin(-1*theta*Math.PI);
};

Transform.prototype.out2in = function(outPt) {
    var x = outPt.x - this.originX;
    var y = outPt.y - this.originY;
    if (this.mirrorY) {
        y = -y;
    }
    x = x/this.scale;
    y = y/(this.scale*this.squashY);
    var rotatedX = this.cosTIn*x - this.sinTIn*y;
    var rotatedY = this.sinTIn*x + this.cosTIn*y;
    var inPt = new Point(rotatedX, rotatedY);
    return inPt;
};
Transform.prototype.in2out = function(inPt) {
    var x = inPt.x;
    var y = inPt.y;
    //var theta = -1*this.theta * Math.PI;
    var rotatedX = this.cosTOut*x - this.sinTOut*y;
    var rotatedY = this.sinTOut*x + this.cosTOut*y;
    //var rotatedX = Math.cos(theta)*x - Math.sin(theta)*y;
    //var rotatedY = Math.sin(theta)*x + Math.cos(theta)*y;
    if (this.mirrorY) {
        rotatedY = -rotatedY;
    }
    rotatedX = rotatedX*this.scale;
    rotatedY = rotatedY*this.scale*this.squashY;
    rotatedX += this.originX;
    rotatedY += this.originY;
    var outPt = new Point(rotatedX, rotatedY);
    return outPt;
};
Transform.prototype.shift = function(dx, dy) {
    if (dx) {
        this.originX += dx;
    }
    if (dy) {
        this.originY += dy;
    }
};
Transform.prototype.setInAtOut = function(inPt, outPt) {
    var curPt = this.in2out(inPt);
    this.shift(outPt.x - curPt.x, outPt.y - curPt.y);
    //this.shift(setPt[0]-curPt[0], setPt[1]-curPt[1]);
};
Transform.prototype.rotateAround = function(theta, outPt) {
    if (!theta) {
        return;
    }
    var inPt;
    if (outPt) {
        inPt = this.out2in(outPt);
    }
    this.setTheta(this.theta + theta);
    //this.theta += theta;
    if (inPt) {
        this.setInAtOut(inPt, outPt);
    }
};
Transform.prototype.scaleAround = function(dScale, outPt) {
    if (!dScale || this.scale+dScale < 0) {
        return;
    }
    var inPt;
    if (outPt) {
        inPt = this.out2in(outPt);
    }
    this.scale += dScale;
    if (inPt) {
        this.setInAtOut(inPt, outPt);
    }
};

muleObj.geometry.Transform = Transform;

})(muleObj);
