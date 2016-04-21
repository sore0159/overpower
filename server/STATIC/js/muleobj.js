//(function() {
    var muleElement = document.getElementById('MULEELEMENT');
    if (!muleElement) {
        console.log("NO MULE ELEMENT FOUND: ABORTING");
        //return;
    }
    if (!muleElement.MULEOBJECT) {
        muleElement.MULEOBJECT = {};
        console.log("CREATING MULE OBJECT");
    }
    var muleObj = muleElement.MULEOBJECT;


//})();
