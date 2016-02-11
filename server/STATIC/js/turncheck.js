(function() {
    var canvas = document.getElementById('mainscreen');
    var blockButton = document.getElementById('blockbutton');
    blockButton.onclick = function() {
        document.location.reload(true);

    };
    canvas.overpowerData = {};
    canvas.overpowerData.game = {turn:curTurn, gid: curGame};
    canvas.turnCheck(true);

    canvas.blockScreen = function(text, buttonT) {
        var screen = document.getElementById('fullscreenblock');
        var blocktext = document.getElementById('blockertext');
        var blockbutton = document.getElementById('blockbutton');
        blocktext.innerHTML = text;
        if (buttonT) {
            blockbutton.textContent = buttonT;
            blockbutton.style.display = "inline";
        } else {
            blockbutton.style.display = "none";
        }
        screen.style.display = 'block';

    };

    canvas.unblockScreen = function() {
        var screen = document.getElementById('fullscreenblock');
        screen.style.display = 'none';
    };
})();
