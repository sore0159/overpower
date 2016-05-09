
(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.ajax) {
    muleObj.ajax = {};
}
if (!muleObj.html) {
    muleObj.html = {};
}


if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}
if (!muleObj.overpower.html) {
    muleObj.overpower.html = {};
}

var html = muleObj.html;
var overpower = muleObj.overpower;

var data = overpower.data;
var ophtml = overpower.html;
var net = overpower.net;




ophtml.infobox.overview = {
    turnbutton: new html.Tree("turnbutton"),
    turntogglebutton: new html.Tree("turntogglebutton"),
    turnconfirmbutton: new html.Tree("turnconfirmbutton"),
    reportsbutton: new html.Tree("reportsbutton"),
    powerbutton: document.getElementById("powerorderbutton"),
    powercancelbutton: new html.Tree("powercancelbutton"),
};

var overview = ophtml.infobox.overview;

ophtml.target2Button(null, overview.powerbutton);
overview.powercancelbutton.setClick(function() {
    if (data.power.planet) {
        overpower.commands.setPowerOrder(data.power.planet, 0);
    }
});

overview.reportsbutton.setClick(function() {
    console.log("VIEW REPORTS GO!");
});

overview.turntogglebutton.setClick(function() {
    var buffer = data.factions.myFaction.donebuffer;
    if (buffer) {
        net.putTurnBuffer(0);
    } else {
        net.putTurnBuffer(1);
    }
});

overview.turntogglebutton.render = function() {
    var buffer = data.factions.myFaction.donebuffer;
    if (buffer) {
        this.elem.textContent = "Set Not Complete";
    } else {
        this.elem.textContent = "Set Complete";
    }
};

overview.turnconfirmbutton.setClick(function() {
    net.putTurnBuffer(overview.turnbutton.userBuffer);
});
overview.turnbutton.setWheel(function(up) {
    if (up > 0) {
        if (!overview.turnbutton.userBuffer) {
            overview.turnbutton.userBuffer = data.factions.myFaction.donebuffer+1;
            if (!overview.turnbutton.userBuffer) {
                overview.turnbutton.userBuffer = 1;
            }
        } else if (overview.turnbutton.userBuffer === -1) {
            overview.turnbutton.userBuffer = 1;
        } else {
            overview.turnbutton.userBuffer += 1;
        }
    } else {
         if (!overview.turnbutton.userBuffer) {
            overview.turnbutton.userBuffer = data.factions.myFaction.donebuffer-1;
        } else if (overview.turnbutton.userBuffer !== -1) {
            overview.turnbutton.userBuffer -= 1;
        }
        if (overview.turnbutton.userBuffer < 1) {
            overview.turnbutton.userBuffer = -1;
        }
    }
    overview.turnbutton.render();
});
overview.turnbutton.render = function() {
    var buffer = data.factions.myFaction.donebuffer;
    var userBuffer = overview.turnbutton.userBuffer;
    if (!buffer) {
        delete overview.turnbutton.userBuffer;
        overview.turnbutton.style.display = 'none';
        overview.turnconfirmbutton.style.display = 'none';
        return;
    }
    if (!userBuffer || buffer === userBuffer) {
        overview.turnconfirmbutton.style.display = 'none';
    } else {
        overview.turnconfirmbutton.style.display = 'inline';
        if (userBuffer === -1) {
            overview.turnconfirmbutton.elem.textContent = "Click to confirm all turns complete";
        } else {
            overview.turnconfirmbutton.elem.textContent = "Click to confirm "+userBuffer+" turns complete";
        }
    }
    if (buffer === 1) {
        overview.turnbutton.elem.textContent = "Create Buffer";
    } else {
        overview.turnbutton.elem.textContent = "Modify Buffer";
    }
    overview.turnbutton.style.display = 'inline';
};

overview.render = function() {
    overview.renderScore();
    overview.renderTurn();
    overview.renderReports();
    overview.renderPowerOrder();
};
overview.renderScore = function() {
    html.setText("scoretext", data.factions.myFaction.score);
    html.setText("towintext", data.game.towin);
    if (data.game.turn > 1) {
        html.setText("highscoretext", data.game.highscore);
    } else if (data.game.turn === 1) {
        html.setText("highscoretext", "-");
    }
};
overview.renderTurn = function() {
    delete overview.turnbutton.elem.userBuffer;
    html.setText("turntext", data.game.turn);
    var buffer = data.factions.myFaction.donebuffer;
    if (buffer === -1) {
        html.setText("turnstatustext", "All Turns Complete");
    } else if (buffer === 1) {
        html.setText("turnstatustext", "Turn Complete");
    } else if (buffer) {
        html.setText("turnstatustext", ""+buffer+" Turns Complete");
    } else {
        html.setText("turnstatustext", "Turn In Progress");
    }
    overview.turntogglebutton.render();
    overview.turnbutton.render();
};
overview.renderReports = function() {
    if (data.game.turn > 1) {
        html.setText("reportnum", data.battleReports.length + data.launchReports.length);
        overview.reportsbutton.style.display = 'inline';

    } else {
        overview.reportsbutton.style.display = 'none';
        html.setText("reportnum", "-");
    }
};
overview.renderPowerOrder = function() {
    var elem = document.getElementById("powerordertext");
    if (data.isPOUseful()) {
        elem.textContent = ((data.power.type === -1) ? "Tachyons ("+data.power.planet.tachyons+")": "Antimatter ("+data.power.planet.antimatter+")");
        elem.className = "";
        overview.powerbutton.style.display = 'inline';
        overview.powercancelbutton.style.display = 'inline';
        overview.powerbutton.setPlanet(data.power.planet);
    } else {
        elem.textContent = "Not yet set";
        elem.className = "alert";
        overview.powercancelbutton.style.display = 'none';
        overview.powerbutton.style.display = 'none';
    }
};


})(muleObj);

