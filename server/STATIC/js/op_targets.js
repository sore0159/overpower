
(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.html) {
    muleObj.html = {};
}

var html = muleObj.html;
var overpower = muleObj.overpower;

if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}
if (!muleObj.overpower.html) {
    muleObj.overpower.html = {};
}
if (!muleObj.overpower.net) {
    muleObj.overpower.net = {};
}


var data = overpower.data;
var ophtml = overpower.html;
var net = overpower.net;

ophtml.infobox.orders = {
    launchbox: new html.Tree("launchorderbox"),
    powerbox: new html.Tree("powerorderbox"),
};

var launchbox = ophtml.infobox.orders.launchbox;
launchbox.render = function() {
    launchbox.clear();
    var tOrd = data.targets.order;
    if (!tOrd) {
        launchbox.style.display = "none";
        return;
    }
    launchbox.style.display = "block";
    launchbox.setText("ORDER:"+JSON.stringify(tOrd));
};
var powerbox = ophtml.infobox.orders.powerbox;
powerbox.render = function() {
    powerbox.clear();
    var pl = data.power.planet;
    if (data.targets.secondary && data.targets.secondary.planet.myControl > 0) {
        pl = data.targets.secondary.planet;
    } else if (!data.isPOUseful() && (!data.targets.secondary && data.targets.primary && data.targets.primary.planet && data.targets.primary.planet.myControl > 0)) {
        pl = data.targets.primary.planet;
    }  else {
        powerbox.style.display = "none";
        return;
    }
    powerbox.style.display = "block";
    powerbox.setText("PLANET:"+JSON.stringify(pl));
};

ophtml.infobox.targets = {
    primary: {
        box: new html.Tree("lefttarget"),
    },
    secondary: {
        box: new html.Tree("righttarget"),
    },
};
var targets = ophtml.infobox.targets;

targets.render = function() {
    targets.primary.render();
    targets.secondary.render();
    launchbox.render();
    powerbox.render();
};
targets.primary.render = function() {
    var box = targets.primary.box;
    box.clear();
    var tDat = data.targets.primary;
    if (!tDat) {
        box.style.visibility = 'hidden';
        return;
    }
    box.style.visibility = 'visible';
    box.spurClass("p", "target1", "Primary Target: ("+tDat.hex.x+","+tDat.hex.y+")");
    if (tDat.dist) {
        box.spur("p", "Distance from secondary target: "+tDat.dist + " sector"+((tDat.dist === 1) ? "":"s" )+" ("+tDat.turns+" turn"+((tDat.turns === 1) ? "":"s")+")");
    }
    if (tDat.planet) {
        box.spurElement(ophtml.jumperButton(tDat.planet));
    }
};
targets.secondary.render = function() {
    var box = targets.secondary.box;
    box.clear();
    var tDat = data.targets.secondary;
    if (!tDat) {
        box.style.visibility = 'hidden';
        return;
    }
    box.style.visibility = 'visible';
    box.spurClass("p", "target2", "Secondary Target: ("+tDat.hex.x+","+tDat.hex.y+")");
    box.spurElement(ophtml.jumperButton(tDat.planet));
    
};

})(muleObj);
