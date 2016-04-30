(function(muleObj) {

var html = muleObj.html;
var geometry = muleObj.geometry;

var stars = {};
muleObj.overpower.stars = stars;

var canvas = document.getElementById('starscreen');
stars.screen = new html.ScreenTransform(canvas, 0,0, 0, 1, 0.55);
stars.locations = [[],[],[],[]];
stars.sizeFilled = 0;

stars.moreStars = function() {
    var starW = (stars.screen.canvas.width > stars.screen.canvas.height) ? stars.screen.canvas.width : stars.screen.canvas.height;
    starW = 2 * starW;
    if (starW <= stars.sizeFilled) {
        return;
    }
    var newRad = starW - stars.sizeFilled;
    var newArea = (starW * starW) - (stars.sizeFilled * stars.sizeFilled);
    var density;
    for (var j = 0; j < 4; j++ ) {
        if (j === 0) {
            density = 0.00625;
        } else if (j === 1) {
            density = 0.003125;
        } else if (j === 2) {
            density = 0.0015625;
        } else if (j === 3) {
            density = 0.00038125;
            //density = 0.00078125;
        }
        density *= 4;
        var numStars = density * newArea;
        for (var i = 0; i< numStars; i++) {
            var x, y;
            x = newRad * 2 * (Math.random() - 0.5);
            y = newRad * 2 * (Math.random() - 0.5);
            if (x > 0) {
                x += stars.sizeFilled;
            } else {
                x -= stars.sizeFilled;
            }
            if (y > 0) {
                y += stars.sizeFilled;
            } else {
                y -= stars.sizeFilled;
            }
            stars.locations[j].push(new geometry.Point(x,y));
        }
    }
    stars.sizeFilled = starW;
};


stars.render = function() {
    var canvas = stars.screen.canvas;
    var ctx = canvas.getContext('2d');
    ctx.fillStyle = "rgb(0,0,0)";
    ctx.fillRect(0,0,canvas.width, canvas.height);
    var img = ctx.getImageData(0,0,canvas.width, canvas.height);
    var bright;
    function drawStar(starPt) {
        var drawPt = stars.screen.transform.in2out(starPt);
        if (drawPt.x < 0 || drawPt.x > canvas.width) {
            return;
        }
        if (drawPt.y < 0 || drawPt.y > canvas.height) {
            return;
        }
        var index = (Math.floor(drawPt.x) + Math.floor(drawPt.y) * canvas.width) * 4;
        var bUse = Math.floor( bright * (0.25 + 0.75*Math.random()) );
        bUse += 10;
        img.data[index] = bUse;
        img.data[index+1] = bUse;
        img.data[index+2] = bUse;
        img.data[index+3] = 255;
    }
    var list = stars.locations[0];
    bright = 25;
    list.forEach(drawStar);
    list = stars.locations[1];
    bright = 50;
    list.forEach(drawStar);
    list = stars.locations[2];
    bright = 100;
    list.forEach(drawStar);
    list = stars.locations[3];
    bright = 200;
    list.forEach(drawStar);
    ctx.putImageData(img, 0, 0);
};

stars.animate = function() {
    var stars = muleObj.overpower.stars;
    stars.screen.rotate(0.00005);
    stars.render();
    if (!stars.stopAnimation) {
        window.requestAnimationFrame(stars.animate);
    }
};

stars.moreStars();
stars.animate();

})(muleObj);
