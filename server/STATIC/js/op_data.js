(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}

if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}

var overpower = muleObj.overpower;
var geometry = muleObj.geometry;
var data = muleObj.overpower.data;

data.targets = {};

data.targets.setT1 = function(hex) {
    hex = data.targets.assist(hex);
    data.targets.primary = { hex: hex };

    var planet = data.planetGrid.getHex(hex);
    if (planet) {
        data.targets.primary.planet = planet;
    }
};
data.targets.setT2 = function(hex, force) {
    hex = data.targets.assist(hex);
    var planet = data.planetGrid.getHex(hex);
    if (!planet) {
        if (force) {
            delete data.targets.secondary;
        }
        return;
    }
    data.targets.secondary = { hex: hex, planet: planet };
};

data.targets.assist = function(hex) {
    if (overpower.map.isShowGrid() || !data.planetGrid || data.planetGrid.getHex(hex)) {
        return hex;
    }
    var search;
    var scale = overpower.map.getScale();
    if (scale > 5) {
        search = 1;
    } else if (scale > 2.5) {
        search = 2;
    } else {
        search = 3;
    }
    for (var i = 1; i <= search; i += 1) {
        var list = hex.ring(i);
        for (var j = 0; j < list.length; j += 1) {
            if (data.planetGrid.getHex(list[j])) {
                return list[j];
            }
        }
    }
    return hex;
};

data.targets.isT1 = function(hex) {
    return (this.primary && hex.eq(this.primary.hex));
};
data.targets.isT2 = function(hex) {
    return (this.secondary && hex.eq(this.secondary.hex));
};

data.parseFullView = function(fullView) {
    //console.log("GOT", JSON.stringify(fullView));
    data.planetGrid = new geometry.HexMap();
    fullView.planetviews.forEach(function(pv) {
        pv.hex = new geometry.Hex(pv.loc[0], pv.loc[1]);
        data.planetGrid.setHex(pv.hex, pv);
    });
    data.planetGrid.forEach(function(pv, hex) {
        //console.log("PV", pv, "AT", hex);
    });

};


})(muleObj);
