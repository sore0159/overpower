
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
    reportsbutton: new html.Tree("reportsbutton"),
    powerbutton: document.getElementById("powerorderbutton"),
};
ophtml.jumperButton(null, ophtml.infobox.overview.powerbutton);

var overview = ophtml.infobox.overview;

//overview.powerbutton.setClick(function() {
    //console.log("VIEW REPORTS GO!");
//});

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

overview.turnbutton.setClick(function() {
    overview.turnbutton.render();
    net.putTurnBuffer(this.userBuffer);
});
overview.turnbutton.setWheel(function(up) {
    if (up > 0) {
        if (!this.userBuffer) {
            this.userBuffer = data.factions.myFaction.donebuffer+1;
            if (!this.userBuffer) {
                this.userBuffer = 1;
            }
        } else if (this.userBuffer === -1) {
            this.userBuffer = 1;
        } else {
            this.userBuffer += 1;
        }
    } else {
         if (!this.userBuffer) {
            this.userBuffer = data.factions.myFaction.donebuffer-1;
        } else if (this.userBuffer !== -1) {
            this.userBuffer -= 1;
        }
        if (!this.userBuffer) {
            this.userBuffer = -1;
        }
    }
    overview.turnbutton.render();
});
overview.turnbutton.render = function() {
    var buffer = data.factions.myFaction.donebuffer;
    var userBuffer = this.elem.userBuffer;
    if (!buffer) {
        this.style.display = 'none';
    } else if (!userBuffer || buffer === userBuffer) {
        this.style.display = 'inline';
        this.elem.className = "inactive";
        if (buffer === 1) {
            this.elem.textContent = "Set Buffer";
        } else {
            this.elem.textContent = "Change Buffer";
        }
    } else {
        this.style.display = 'block';
        this.elem.className = "active";
        if (userBuffer === -1) {
            this.elem.textContent = "Click to confirm all turns complete";
        } else {
            this.elem.textContent = "Click to confirm "+userBuffer+" turns complete";
        }
    }
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
    overview.turnbutton.render();
    overview.turntogglebutton.render();
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
        elem.textContent = ((data.power.type === -1) ? "Tachyons": "Antimatter");
        elem.className = "";
        overview.powerbutton.style.display = 'inline';
        overview.powerbutton.setPlanet(data.power.planet);
    } else {
        elem.textContent = "Not yet set";
        elem.className = "alert";
        overview.powerbutton.style.display = 'none';
    }
};


})(muleObj);

