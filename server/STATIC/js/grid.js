/*
        ____
       /\  /\
  ____/__\/__\
 /\  /\  /\  /
/__\/__\/__\/
\  /\  / 
 \/__\/

*/

(function() {
function hexSteps(hex1, hex2) {
    var z1  = -hex1[0] - hex1[1];
    var z2  = -hex2[0] - hex2[1];
    var dx = Math.abs(hex2[0] - hex1[0]);
    var dy = Math.abs(hex2[1] - hex1[1]);
    var dz = Math.abs(z2 - z1);
    return (dx + dy + dz)/2;
}

function h2Center(hex) {
    var px = 1.5 * hex[0];
    // 1^2 = (.5)^2 + (dy^2)
    // dy = sqrt( .75 ) =  .86602540378
    var py = (2*hex[1] + hex[0])*0.86602540378;
    return [px, py];
}

function h2Corners(hex) {
    var hexTriH = 0.86602540378;
    var center = h2Center(hex);
    var pts = [
        [center[0]+1, center[1]],
        [center[0]+0.5, center[1]+hexTriH],
        [center[0]-0.5, center[1]+hexTriH],
        [center[0]-1, center[1]],
        [center[0]-0.5, center[1]-hexTriH],
        [center[0]+0.5, center[1]-hexTriH],
    ];
    return pts;
}

function p2Hex(pt) {
	//       __
	//      | \|  width =  1.5
	//box = |_/|  height = sqrt(3) = 2*hexTriH;
    //
    var hexTriH = 0.86602540378;
    var x = pt[0] + 0.5;
    var y = pt[1] + hexTriH;
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
    return [hx, hy];
}

function Grid(scale, originX, originY, theta, squashY, noMirrorY) {
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

Grid.prototype.out2in = function(pt) {
    var x = pt[0] - this.originX;
    var y = pt[1] - this.originY;
    if (this.mirrorY) {
        y = -y;
    }
    x = x/this.scale;
    y = y/(this.scale*this.squashY);
    var theta = this.theta * Math.PI;
    var rotatedX = Math.cos(theta)*x - Math.sin(theta)*y;
    var rotatedY = Math.sin(theta)*x + Math.cos(theta)*y;
    return [rotatedX, rotatedY];
};
Grid.prototype.in2out = function(pt) {
    var x = pt[0];
    var y = pt[1];
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
    return [rotatedX, rotatedY];
};
Grid.prototype.setHexAt = function(hex, setPt) {
    var curPt = this.h2Center(hex);
    this.shift(setPt[0]-curPt[0], setPt[1]-curPt[1]);
};
Grid.prototype.setInPtAt = function(inPt, setPt) {
    var curPt = this.in2out(inPt);
    this.shift(setPt[0]-curPt[0], setPt[1]-curPt[1]);
};
Grid.prototype.setOutPtAt = function(outPt, setPt) {
    this.shift(setPt[0]-outPt[0], setPt[1]-outPt[1]);
};
Grid.prototype.p2Hex = function(pt) {
    var inPt = this.out2in(pt);
    return p2Hex(inPt);
};
Grid.prototype.h2Center = function(hex) {
    var inPt = h2Center(hex);
    return this.in2out(inPt);
};
Grid.prototype.h2Corners = function(hex) {
    var inList = h2Corners(hex);
    var outList = [];
    var grid = this;
    inList.forEach(function(inPt) {
        var outPt = grid.in2out(inPt);
        outList.push(outPt);
    });
    return outList;
};
Grid.prototype.shift = function(dx, dy) {
    if (dx) {
        this.originX += dx;
    }
    if (dy) {
        this.originY += dy;
    }
};
Grid.prototype.rotateAround = function(theta, aboutPt) {
    if (!theta) {
        return;
    }
    var curIn;
    if (aboutPt) {
        curIn = this.out2in(aboutPt);
    }
    this.theta += theta;
    if (curIn) {
        this.setInPtAt(curIn, aboutPt);
    }
};
Grid.prototype.scaleAround = function(dScale, aboutPt) {
    if (!dScale || this.scale+dScale < 0) {
        return;
    }
    var curIn;
    if (aboutPt) {
        curIn = this.out2in(aboutPt);
    }
    this.scale += dScale;
    if (curIn) {
        this.setInPtAt(curIn, aboutPt);
    }
};
Grid.prototype.visibleHexList = function(minX, minY, maxX, maxY) {
    var minObj = {};
    var maxObj = {};
    var grid = this;
    function xChecker(x, y) {
        var hex = grid.p2Hex([x,y]);
        var key = ""+hex[0];
        if (!minObj[key] || minObj[key][1] > hex[1]) {
            minObj[key] = hex;
        }
        if (!maxObj[key] || maxObj[key][1] < hex[1]) {
            maxObj[key] = hex;
        }
    }
    var y = 0;
    var x = minX;
    for (y = minY; y <=maxY; y+=1) {
        xChecker(x,y);
    }
    x = maxX;
    for (y = minY; y <=maxY; y+=1) {
        xChecker(x,y);
    }
    y = minY;
    for (x = minX; x <=maxX; x+=1) {
        xChecker(x,y);
    }
    y = maxY;
    for (x = minX; x <=maxX; x+=1) {
        xChecker(x,y);
    }
    var hexList = [];
    Object.keys(minObj).forEach(function(key) {
        var minHex = minObj[key];
        var maxHex = maxObj[key];
        for (i = minHex[1]; i<=maxHex[1];i++) {
            var hex = [minHex[0], i];
            hexList.push(hex);
        }
    });
    return hexList;
};
Grid.prototype.stepsBetween = hexSteps;
Grid.prototype.ptsEq = function(pt1, pt2) {
    if (!pt1 || !pt2) {
        return false;
    }
    return (pt1[0] === pt2[0] && pt1[1] === pt2[1]);
};

var canvas = document.getElementById('mainscreen');
var grid = new Grid(40, canvas.width/2, canvas.height/2, 0.25, 0.55);
canvas.muleGrid = grid;

canvas.drawGrid = function(hexList) {
    var ctx = this.getContext('2d');
    ctx.clearRect(0,0,this.width, this.height);
    var scale = this.muleGrid.scale;
    if (scale < 25) {
        ctx.lineWidth = 0.5;
    } else if (scale < 35) {
        ctx.lineWidth = 1;
    } else {
        ctx.lineWidth = 1.5;
    }
    var grid = this.muleGrid;
    var path = new Path2D();
    function drawHex(hex) {
        var corners = grid.h2Corners(hex);
        path.moveTo(corners[5][0], corners[5][1]);
        corners.forEach(function(pt) {
            path.lineTo(pt[0], pt[1]);
        });
    }
    hexList.forEach(drawHex);

    ctx.strokeStyle = "#2f2f8f";
    ctx.stroke(path);
};

canvas.hexPath = function(hex) {
    var path = new Path2D();
    var corners = this.muleGrid.h2Corners(hex);
    path.moveTo(corners[5][0], corners[5][1]);
    corners.forEach(function(pt) {
        path.lineTo(pt[0], pt[1]);
    });
    return path;
};

canvas.centerHex = function(hex) {
    canvas.muleGrid.setHexAt(hex, [this.width/2, this.height/2]);
};
canvas.centerPt = function(pt) {
    canvas.muleGrid.setOutPtAt(pt, [this.width/2, this.height/2]);
};


})();
