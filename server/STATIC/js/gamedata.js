(function() {

var canvas = document.getElementById('mainscreen');
var data = canvas.overpowerData;

data.setTargetOne = function(hex, help) {
    var grid = canvas.muleGrid;
    if (!help) {
        this.targetOne = hex;
        return;
    }
    var tolerance;
    if (grid.scale > 2.5) {
        tolerance = 2;
    } else if (grid.scale > 1.75) {
        tolerance = 3;
    } else {
        tolerance = 4;
    }
    var exact, closest, dist;
    data.planetviews.forEach(function(planet) {
        if (exact) {
            return;
        }
        if (grid.ptsEq(planet.loc, hex)) {
            exact = hex;
            return;
        }
        var steps = grid.stepsBetween(planet.loc, hex);
        if (!dist || steps < dist) {
            closest = planet.loc;
            dist = steps;
        }
    });
    if (exact) {
        this.targetOne = exact;
        return;
    }
    data.shipviews.forEach(function(ship) {
        if (exact) {
            return;
        }
        if (ship.loc) {
            if (grid.ptsEq(ship.loc, hex)) {
                exact = hex;
                return;
            }
            var steps = grid.stepsBetween(ship.loc, hex);
            if (!dist || steps < dist) {
                closest = ship.loc;
                dist = steps;
            }
        }
    });
    if (exact) {
        this.targetOne = exact;
        return;
    }

    if (dist < tolerance) {
        this.targetOne = closest;
        return;
    }
    this.targetOne = hex;
};

data.setTargetTwo = function(hex, help, drop) {
    var grid = canvas.muleGrid;
    var tolerance;
    if (grid.scale > 2.5) {
        tolerance = 2;
    } else if (grid.scale > 1.75) {
        tolerance = 3;
    } else {
        tolerance = 4;
    }
    var exact, closest, dist;
    data.planetviews.forEach(function(planet) {
        if (exact) {
            return;
        }
        if (grid.ptsEq(planet.loc, hex)) {
            exact = hex;
            return;
        }
        var steps = grid.stepsBetween(planet.loc, hex);
        if (!dist || steps < dist) {
            closest = planet.loc;
            dist = steps;
        }
    });
    if (exact) {
        this.targetTwo = exact;
        return;
    }
    if (help && dist < tolerance) {
        this.targetTwo = closest;
        return;
    }
    if (drop) {
        this.targetTwo = null;
        return;
    }
};


canvas.parseOPData = function() {
    var data = this.overpowerData;
    console.log("DATA:", data);
};





canvas.parseOPData();


})();

