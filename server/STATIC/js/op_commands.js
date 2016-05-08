
(function(muleObj) {

var html = muleObj.html;
var geometry = muleObj.geometry;
var overpower = muleObj.overpower;

if (!overpower.commands) {
    overpower.commands = {};
}
var commands = overpower.commands;

commands.setT1 = function(hex) {
    overpower.data.targets.setT1(hex);
    overpower.html.infobox.targets.render();
    overpower.map.redraw = true;
};
commands.setT2 = function(hex, shift) {
    overpower.data.targets.setT2(hex, shift);
    overpower.html.infobox.targets.render();
    overpower.map.redraw = true;
};
commands.swapTargets = function() {
    if (!overpower.data.targets.primary || !overpower.data.targets.primary.planet || !overpower.data.targets.secondary || overpower.data.targets.secondary.copy) {
        return;
    }
    var h1 = overpower.data.targets.primary.hex;
    var h2 = overpower.data.targets.secondary.hex;
    overpower.data.targets.setT1(h2);
    overpower.data.targets.setT2(h1);
    overpower.html.infobox.targets.render();
    overpower.map.redraw = true;
};

commands.modLaunchOrder = function(plus) {
    var delta = (plus) ? 1: -1;
    if (overpower.data.targets.modOrder(delta)) {
        overpower.map.redraw = true;
        overpower.html.infobox.targets.render();
        //overpower.html.infobox.orders.launchbox.render();
    }
};
commands.confirmLaunchOrder = function() {
    if (overpower.data.targets.order && overpower.data.targets.order.modified) {
        overpower.net.putLaunchOrder(overpower.data.targets.order);
    }
};

commands.clickHex = function(hex, button, shift) {
    if (button === 0) { 
        if (shift) {
            commands.setMapCenter(hex);
        } else {
            commands.setT1(hex);
        }
    } else if (button === 2) {
        commands.setT2(hex, shift);
    } else if (button === 1) {
        commands.setMapCenter(hex);
    }
};
commands.setMapCenter = function(hex) {
    overpower.map.center = hex;
};
commands.modMapFrame = function(plus) {
    var delta = (plus) ? 1: -1;
    overpower.map.setFrame(overpower.map.screen.canvas.width + (12*delta), overpower.map.screen.canvas.height + (9* delta));
    overpower.map.redraw = true;
};

commands.rotateMap = function(clockwise) {
    var delta = (clockwise) ? 1: -1;
    delta *= 0.015;
    overpower.map.screen.rotate(delta);
    overpower.map.redraw = true;
    overpower.stars.screen.rotate(delta*0.5);
};
   
commands.zoomMap = function(zoomIn) {
    if (overpower.map.scaleStep(zoomIn)) {
        overpower.map.redraw = true;
    }
};
commands.setPowerOrder = function(planet, power) {
    if (!planet.myControl || planet.myPower === power) {
        return;
    }
    overpower.net.putPowerOrder({hex: planet.hex, type: power});
};

})(muleObj);
