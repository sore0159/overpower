
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
    var elem, text;
    launchbox.clear();
    var tOrd = data.targets.order;
    if (!tOrd) {
        launchbox.style.display = "none";
        return;
    }
    launchbox.style.display = "block";
    if (tOrd.curO) {
        launchbox.addText("Current launch toward ");
    } else {
        launchbox.addText("Potential launch toward ");
    }
    launchbox.spurElement(ophtml.jumperButton(tOrd.targetPL));
    launchbox.addText(" from ");
    launchbox.spurElement(ophtml.target2Button(tOrd.sourcePL)).and("br");
    if (tOrd.curO) {
        launchbox.addDisplay("Current Size", tOrd.curO.size);
    }
    launchbox.addDisplay("Maximum Available", tOrd.size + tOrd.sourcePL.avail);
    launchbox.spur("br");
    text = (!tOrd.curO) ? "Create Launch": "Modify Launch";
    elem = launchbox.spurClass("button", "inactive", text);
    elem.elem.title = "Scroll the mousewheel over this button to change the size of the launch.  Set to 0 to cancel the launch.";
    elem.setWheel(function(up) {
        overpower.commands.modLaunchOrder(up > 0);
    });
    if (tOrd.modified) {
        launchbox.addDisplay(" New Size", tOrd.size, "alert");
        text = (!tOrd.curO) ? "Confirm Order": "Confirm Changes";
        elem = launchbox.spurClass("button", "active", text);
        elem.setClick(overpower.commands.confirmLaunchOrder);
        elem.elem.title = "You must click here to confirm your launch order before it is final";
    }
};
var powerbox = ophtml.infobox.orders.powerbox;
powerbox.render = function() {
    var elem;
    powerbox.clear();
    var pl = data.power.planet;
    var useful = data.isPOUseful();
    if (!useful && data.targets.secondary && data.targets.secondary.planet.myControl > 0) {
        pl = data.targets.secondary.planet;
    } else if (!useful && (!data.targets.secondary && data.targets.primary && data.targets.primary.planet && data.targets.primary.planet.myControl > 0)) {
        pl = data.targets.primary.planet;
    }  else {
        powerbox.style.display = "none";
        return;
    }
    powerbox.style.display = "block";
    powerbox.spurElement(ophtml.target2Button(pl));
    var b1, b2;
    if (pl.myPower === -1) {
        powerbox.addDisplay(" Current Power", "Tachyons ("+pl.tachyons+")");
        powerbox.spur("br");
        b1 = powerbox.spur("button", "Attune Antimatter ("+pl.antimatter+")");
    } else if (pl.myPower === 1) {
        powerbox.addDisplay(" Current Power", "Antimatter ("+pl.antimatter+")");
        powerbox.spur("br");
        b2 = powerbox.spur("button", "Attune Tachyons ("+pl.tachyons+")");
    } else  {
        powerbox.addDisplay(" Current Power", "None"); 
        powerbox.spur("br");
        b1 = powerbox.spur("button", "Attune Antimatter ("+pl.antimatter+")");
        powerbox.addText(" ");
        b2 = powerbox.spur("button", "Attune Tachyons ("+pl.tachyons+")");
    }
    if (b1) {
        b1.setClick(function() {
            overpower.commands.setPowerOrder(pl, 1);
        });
        b1.elem.title = "This planet will not become attuned to this resource until your next turn.  You can only attune one planet per turn.  Planets can only be attuned to one resource: this will override any existing attunement.";
    }
    if (b2) {
        b2.setClick(function() {
            overpower.commands.setPowerOrder(pl, -1);
        });
        b2.elem.title = "This planet will not become attuned to this resource until your next turn.  You can only attune one planet per turn.  Planets can only be attuned to one resource: this will override any existing attunement.";
    }
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
    if (tDat.planet) {
        box.spurElement(ophtml.jumperButton(tDat.planet));
    }
    if (tDat.dist) {
        box.spur("p", "Distance from secondary target: "+tDat.dist + " sector"+((tDat.dist === 1) ? "":"s" )+" ("+tDat.turns+" turn"+((tDat.turns === 1) ? "":"s")+")");
    }
    if (data.targets.primary && data.targets.primary.planet && data.targets.secondary && !data.targets.secondary.copy) {
        var elem = box.spurClass("button", "target2", "Swap Primary/Secondary");
        elem.setClick(overpower.commands.swapTargets);
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
    box.spurElement(ophtml.target2Button(tDat.planet));
    
};

})(muleObj);
