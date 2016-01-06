var curTurn = -1;

function getTurn() {
    var url = "http://1";
    var req = new XMLHttpRequest();
    //req.addEventListener();
    req.open("GET", url, true);
    req.send();
    var turn = 1;
    return turn;
}
