/*$(document).ready(function () {
	'use strict';
	paper.install(window);
	paper.setup(document.getElementById('mainCanvas'));

	var tool = new Tool();
	var c = Shape.Circle(200, 200, 80);
	c.fillColor = 'black';
	var text = new PointText(200, 200)
	text.justification = 'center';
	text.fillColor = 'white';
	text.fontSize = 20;
	text.content = 'hello world';

	tool.onMouseDown = function(event) {
		var c = Shape.Circle(event.point, 30);
		c.fillColor = 'green';
	};

	paper.view.draw();
});*/

'use strict';
// es6 feature: block-scoped "let" declaration

const sentences = [{ subject: 'LOL', verb: 'is', object: 'great' }, { subject: 'Elephants', verb: 'are', object: 'large' }];
// es6 feature: object destructuring
function say({ subject, verb, object }) {
	// es6 feature: template strings
	console.log(`${ subject } ${ verb } ${ object }`);
}
// es6 feature: for..of
for (let s of sentences) {
	say(s);
}