(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.geometry) {
    console.log("OP STARTS FAILED: REQUIRES GEOMETRY");
    return;
}

var stars = {};
muleObj.overpower.starfield = stars;

stars.canvas = document.getElementById('starscreen');
stars.locations = [[],[],[],[]];
for (var j = 0; j < 4; j++ ) {
    var numStars;
    if (j === 0) {
        numStars = 12000;
    } else if (j === 1) {
        numStars = 6000;
    } else if (j === 2) {
        numStars = 3000;
    } else if (j === 3) {
        numStars = 1500;
    }
    for (var i = 0; i< numStars; i++) {
        x = stars.canvas.width*(1 - 2*Math.random());
        y = stars.canvas.height*(1 - 2*Math.random());
        stars.locations[j].push(new muleObj.geometry.Point(x,y));
    }
}
stars.transform = new muleObj.geometry.Transform(1, stars.canvas.width/2, stars.canvas.height/2, 0, 0.55);
stars.draw = function() {
    var canvas = stars.canvas;
    var ctx = canvas.getContext('2d');
    ctx.fillStyle = "rgb(0,0,0)";
    ctx.fillRect(0,0,canvas.width, canvas.height);
    //
    var path = new Path2D();
    function drawStar(starPt) {
        var drawPt = stars.transform.in2out(starPt);
        if (drawPt.x < 0 || drawPt.x > canvas.width) {
            return;
        }
        if (drawPt.y < 0 || drawPt.y > canvas.height) {
            return;
        }
        path.moveTo(drawPt.x, drawPt.y);
        path.arc(drawPt.x + 0.5, drawPt.y, 0.5, 0, 2*Math.PI);
    }
    var list = stars.locations[0];
    ctx.fillStyle = "rgb(25,25,25)";
    list.forEach(drawStar);
    ctx.fill(path);

    path = new Path2D();
    list = stars.locations[1];
    ctx.fillStyle = "rgb(50,50,50)";
    list.forEach(drawStar);
    ctx.fill(path);

    path = new Path2D();
    list = stars.locations[2];
    ctx.fillStyle = "rgb(100,100,100)";
    list.forEach(drawStar);
    ctx.fill(path);

    path = new Path2D();
    list = stars.locations[3];
    ctx.fillStyle = "rgb(200,200,200)";
    list.forEach(drawStar);
    ctx.fill(path);
};

stars.rotate = function(clockwise, scale) {
    var dTheta = (clockwise) ? 1: -1;
    scale = (scale) ? scale : 1;
    dTheta *= scale * 0.01;
    var center = new muleObj.geometry.Point(stars.canvas.width/2, stars.canvas.height/2);
    stars.transform.rotateAround(dTheta, center);
};

stars.draw();

function animateStars() {
    stars.rotate(true, 0.05);
    stars.draw();
    window.requestAnimationFrame(animateStars);
}
//animateStars();

})(muleObj);

