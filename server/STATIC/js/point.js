(function(muleObj) {

if (!muleObj.geometry) {
    muleObj.geometry = {};
}

function Point(x, y) {
    if(x) {
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
Point.prototype.add = function(x, y) {
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
    theta = Math.pi * theta;
    var dx = r * Math.cos(theta);
    var dy = r * Math.sin(theta);
    var pt = new Point(this.x+dx, this.y+dy);
    return pt;
};
Point.prototype.polarTo = function(point) {
    if (this.eq(point)) {
        return [0, 0];
    }
    var r = this.dist(point);
    var cosT = (point.x - this.x) / r;
    var theta = Math.acos(cosT)/ Math.pi;
    if (point.y < this.y) {
        theta = 2 - theta;
    }
    return [r, theta];
};


muleObj.geometry.Point = Point;


function GridTransform(scale, originX, originY, theta, squashY, noMirrorY) {
    if (scale) {
        this.scale = scale;
    } else {
        this.scale = 1;
    }
    if (originX) {
        this.originX = originX;
    } else {
        this.originX = 0;
    }
    if (originY) {
        this.originY = originY;
    } else {
        this.originY = 0;
    }
    if (theta) {
        this.theta = theta;
    } else {
        this.theta = 0;
    }
    if (noMirrorY) {
        this.mirrorY = false;
    } else {
        this.mirrorY = true;
    }
    if (squashY) {
        this.squashY = squashY;
    } else {
        this.squashY = 1;
    }
}

GridTransform.prototype.out2in = function(outPt) {
    var x = outPt.x - this.originX;
    var y = outPt.y - this.originY;
    if (this.mirrorY) {
        y = -y;
    }
    x = x/this.scale;
    y = y/(this.scale*this.squashY);
    var theta = this.theta * Math.PI;
    var rotatedX = Math.cos(theta)*x - Math.sin(theta)*y;
    var rotatedY = Math.sin(theta)*x + Math.cos(theta)*y;
    var inPt = new Point(rotatedX, rotatedY);
    return inPt;
};
GridTransform.prototype.in2out = function(inPt) {
    var x = inPt.x;
    var y = inPt.y;
    var theta = -1*this.theta * Math.PI;
    var rotatedX = Math.cos(theta)*x - Math.sin(theta)*y;
    var rotatedY = Math.sin(theta)*x + Math.cos(theta)*y;
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
GridTransform.prototype.shift = function(dx, dy) {
    if (dx) {
        this.originX += dx;
    }
    if (dy) {
        this.originY += dy;
    }
};
GridTransform.prototype.setInAtOut = function(inPt, outPt) {
    var curPt = this.in2out(inPt);
    this.shift(outPt.x - curPt.x, outPt.y - curPt.y);
    //this.shift(setPt[0]-curPt[0], setPt[1]-curPt[1]);
};
GridTransform.prototype.rotateAround = function(theta, outPt) {
    if (!theta) {
        return;
    }
    var inPt;
    if (outPt) {
        inPt = this.out2in(outPt);
    }
    this.theta += theta;
    if (inPt) {
        this.setInAtOut(inPt, outPt);
    }
};
GridTransform.prototype.scaleAround = function(dScale, outPt) {
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

muleObj.geometry.GridTransform = GridTransform;

})(muleObj);
