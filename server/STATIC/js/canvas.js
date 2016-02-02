(function() {
    // ----------- SETUP ------------- //
    var canvas = document.getElementById('mainscreen');
    var ctx = canvas.getContext('2d');
    ctx.fillStyle = "rgb(0,0,0)";
    ctx.fillRect(0,0,canvas.width, canvas.height);
    for (var j = 0; j < 4; j++ ) {
        var path = new Path2D();
        var numStars;
        if (j === 0) {
            ctx.fillStyle = "rgb(25,25,25)";
            numStars = 6000;
        } else if (j === 1) {
            ctx.fillStyle = "rgb(50,50,50)";
            numStars = 3000;
        } else if (j === 2) {
            ctx.fillStyle = "rgb(100,100,100)";
            numStars = 1500;
        } else if (j === 3) {
            ctx.fillStyle = "rgb(200,200,200)";
            numStars = 750;
        }
        for (var i = 0; i< numStars; i++) {
            var width = 0.5;
            x = canvas.width* Math.random();
            y = canvas.height* Math.random();
            path.moveTo(x, y);
            path.arc(x+width, y, width, 0, 2*Math.PI);
        }
        ctx.fill(path);
    }
    function mapClick(event) {
        var clickx = event.pageX - canvas.offsetLeft;
        var clicky = event.pageY - canvas.offsetTop;
        console.log("click", event.button, clickx, clicky);
    }
    function mapWheel(event) {
        event.preventDefault();
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
        console.log("wheel", up);
        return false;
    }
    canvas.onmousedown = mapClick;
    canvas.oncontextmenu= function() { return false; };
    canvas.onmousewheel = mapWheel;
    canvas.onDOMMouseScroll = mapWheel;
    canvas.addEventListener("DOMMouseScroll", mapWheel);
})();
