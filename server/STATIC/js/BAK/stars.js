(function() {
    // ----------- SETUP ------------- //
    var canvas = document.getElementById('starscreen');
    canvas.drawStars = function() {
        var ctx = this.getContext('2d');
        ctx.fillStyle = "rgb(0,0,0)";
        ctx.fillRect(0,0,this.width, this.height);
        //
        var path = new Path2D();
        function drawStar(pt) {
            path.moveTo(pt[0], pt[1]);
            path.arc(pt[0]+0.5, pt[1], 0.5, 0, 2*Math.PI);
        }
        var list = this.starPoints[0];
        ctx.fillStyle = "rgb(25,25,25)";
        list.forEach(drawStar);
        ctx.fill(path);

        path = new Path2D();
        list = this.starPoints[1];
        ctx.fillStyle = "rgb(50,50,50)";
        list.forEach(drawStar);
        ctx.fill(path);

        path = new Path2D();
        list = this.starPoints[2];
        ctx.fillStyle = "rgb(100,100,100)";
        list.forEach(drawStar);
        ctx.fill(path);

        path = new Path2D();
        list = this.starPoints[3];
        ctx.fillStyle = "rgb(200,200,200)";
        list.forEach(drawStar);
        ctx.fill(path);
    };
    canvas.starPoints = [[],[],[],[]];
    for (var j = 0; j < 4; j++ ) {
        var numStars;
        if (j === 0) {
            numStars = 6000;
        } else if (j === 1) {
            numStars = 3000;
        } else if (j === 2) {
            numStars = 1500;
        } else if (j === 3) {
            numStars = 750;
        }
        for (var i = 0; i< numStars; i++) {
            x = canvas.width* Math.random();
            y = canvas.height* Math.random();
            canvas.starPoints[j].push([x,y]);
        }
    }
    canvas.drawStars();
})();
