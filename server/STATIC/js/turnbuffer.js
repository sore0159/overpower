(function() {

var canvas = document.getElementById('mainscreen');

var boxTurnBuffer = document.getElementById('bufferbox');
boxTurnBuffer.redraw = function() {
    var bufferVal = canvas.overpowerData.faction.donebuffer;
    var turnCompText = document.getElementById('turncompletetext');
    if (!bufferVal) {
        this.style.display = "none";
        turnCompText.textContent = "Turn in progress";
        boxTurnButton.textContent = "Set turn complete";
        return;
    } 

    this.style.display = "inline";
    turnBufferButton.clickBuffer = bufferVal;
    turnBufferButton.redraw();
    var bufferText = document.getElementById('buffertext');
    if (bufferVal === 1) {
        boxTurnButton.textContent = "Cancel turn complete";
        turnCompText.textContent = "Turn complete";
        bufferText.textContent = "Set further turns to auto-complete?";
    } else if (bufferVal === -1) {
        boxTurnButton.textContent = "Clear turn completion buffer";
        turnCompText.textContent = "Turns set to auto-complete indefinitely";
        bufferText.textContent = "Set further turns to auto-complete?";
    } else {
        boxTurnButton.textContent = "Clear turn completion buffer";
        turnCompText.textContent = data.faction.donebuffer+" turns set to complete";
    }
    
};


var boxTurnButton = document.getElementById('turnchange');
boxTurnButton.addEventListener("mouseup", turnChange);

function turnChange() {
    var faction = canvas.overpowerData.faction;
    if (!faction.donebuffer) {
        faction.donebuffer = 1;
        faction.done = true;
    } else {
        faction.donebuffer = 0;
        faction.done = false;
    }
    boxTurnBuffer.redraw();
    if (faction.donebuffer) {
        canvas.blockScreen("Checking for new turn...");
    }
    canvas.putJSON("/overpower/json/factions", faction, canvas.putErr, function() {
        if (faction.donebuffer) {
            canvas.turnCheck();
        }
    });
}

var turnBufferButton = document.getElementById('turnbufferbutton');
turnBufferButton.addEventListener("mouseup", bufferClick);
turnBufferButton.addEventListener("DOMMouseScroll", bufferScroll);
turnBufferButton.onmousewheel = bufferScroll;
turnBufferButton.redraw = function() {
    var buffer =  canvas.overpowerData.faction.donebuffer;
    if (this.clickBuffer === buffer) {
        this.textContent = "Mousescroll to set buffer";
        this.style.fontWeight = "normal";
    } else {
        this.textContent = "Click to set buffer to "+this.clickBuffer;
        this.style.fontWeight = "bold";
    }
};

function bufferClick(event) {
    var buffer =  canvas.overpowerData.faction.donebuffer;
    if (!buffer || buffer === this.clickBuffer) {
        return;
    }
    console.log("BUFFER CLICK", this.clickBuffer);
}

function bufferScroll(event) {
    event.preventDefault();
    var gameturn = canvas.overpowerData.game.turn;
    var up;
    if (event.detail) {
        up = -1*event.detail/3;
    } else {
        up = (event.wheelDelta)/120;
    }
    if (up > 0 && up < 1) {
        up = 1;
    } else if (up < 0 && up > -1) {
        up = -1;
    } else if (up === 0) {
        return true;
    }
    if (up === 1) {
        if (this.clickBuffer === -1) {
            this.clickBuffer = 1;
        } else {
            this.clickBuffer += 1;
        }
    } else if (up === -1) {
       if (this.clickBuffer > 1) {
            this.clickBuffer -= 1;
       } else if (this.clickBuffer === 1) {
            this.clickBuffer = -1;
       }
    } else {
        return false;
    }
    this.redraw();
    return false;
}
    
})();
