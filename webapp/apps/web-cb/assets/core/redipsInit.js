seq = ko.observable(0)
/*jslint white: true, browser: true, undef: true, nomen: true, eqeqeq: true, plusplus: false, bitwise: true, regexp: true, strict: true, newcap: true, immed: true, maxerr: 14 */
/*global window: false, REDIPS: true */

/* enable strict mode */
//"use strict";
// define redips_init variable
var redipsObjectSelected;
var redipsInit;
var rdDraggedItem = ko.observable();
// redips initialization
redipsInit = function () {
	// reference to the REDIPS.drag library and message line
	var	rd = REDIPS.drag;
	// how to display disabled elements
	rd.style.borderDisabled = 'solid';	// border style for disabled element will not be changed (default is dotted)
	rd.style.opacityDisabled = 60;		// disabled elements will have opacity effect
	// initialization
	rd.init();
	// only "smile" can be placed to the marked cell
	// rd.mark.exception.d8 = 'smile';
	// prepare handlers

	// rd.event.clicked = function () {
	// 	console.log('Clicked');
	// };

	// rd.event.moved  = function () {
	// 	console.log('Moved');
	// };

	rd.event.dblClicked = function () {
		var str = $(rd.obj).attr("id");
		var code = str.substring(4);
		var res = str.slice(0,4); 
		if (res != "task"){
			Initiative.Get($(rd.obj).attr("id"));
		} else {
			Task.Get(code);
		}
	};

	rd.event.dropped = function () {
		pos = rd.getPosition();
		// console.log(pos);
		rd.posOldRow = pos[4]
		rd.posOldCel = pos[5]
		rd.posNewRow = pos[1]
		rd.posNewCel = pos[2]
		redipsObjectSelected = rd;
		// console.log(rd);
		if(rd.posOldRow != rd.posNewRow || rd.posOldCel != rd.posNewCel){			
			var scID = $(rd.td.target).attr("SCCategory");
			var bdID = $(rd.td.target).attr("BDId");
			var lcID = $(rd.td.target).children().attr("lcID");
			// console.log(bdID, lcID, $(rd.td.target))
			rdDraggedItem({scID:scID,lcID:lcID,lcID:lcID})
			rd.moveObject({
		        obj: rd.obj,
		        target: [0, rd.posOldRow, rd.posOldCel]
		    });

		    $('#mdlcloneormove').modal('show')
		}

	};
};

function rdMove(){
	rd = redipsObjectSelected;
	rd.moveObject({
        obj: rd.obj,
        target: [0, rd.posNewRow, rd.posNewCel]
    });
	var obj = $(rd.obj);
	var target = $(rd.td.current)
	var str = obj.attr("id");
	var code = str.substring(4);
	var initiativeId = "";	
	var filter = "";
	var res = str.slice(0,4); 
	var url = "";
	// console.log("ini target", target)
	if (res != "task"){
		initiativeId = obj.attr("InitiativeID");
		filter = "$.InitiativeID == '" + initiativeId +"'";
		var selected = Enumerable.From(c.AllInitiateSource()).FirstOrDefault(undefined, filter)
		selected.LifeCycleId = target.attr("LCId");
		if (target.parent().attr("SCCategory") != 'All'){
			selected.SCCategory = target.parent().attr("SCCategory");
			selected.BusinessDriverId = target.parent().attr("BDId");
		}
		// selected.BusinessDriverId = target.attr("BDId");
		url = "/web-cb/dashboard/moveupdate";

	} else {
		filter = "$.Id == '" + code +"'";
		var selected = Enumerable.From(c.DataSource().Data.TaskList).FirstOrDefault(undefined, filter)
		selected.LifeCycleId = target.attr("LCId");
		if (target.parent().attr("SCCategory") != 'All'){
			selected.SCCategory = target.parent().attr("SCCategory");
			selected.BusinessDriverId = target.parent().attr("BDId");
		}
		// selected.BusinessDriverId = target.attr("BDId");
		url = "/web-cb/task/moveupdate";
	}

		obj.css("border-color", target.attr("colorcode"));	
	    $('#mdlcloneormove').modal('hide')


	  // console.log(selected, target.attr("LCId"), target.attr("BDId"), target)
	
		ajaxPost(url, selected, 
			function(res){
				if(res.Res != "OK"){
					swal("Move", "Error!", "error");
					return;
				}
				c.GetData();

			}
		)
}

function rdClone(){
	rd = redipsObjectSelected;
	var obj = $(rd.obj);
	console.log(obj)
	var target = $(rd.td.current)
	var str = obj.attr("id");
	var code = str.substring(4);
	var initiativeId =0;	
	var filter = "";
	var res = str.slice(0,4); 
	var url = "";
	var d = rdDraggedItem();
	// obj.css("border-color", target.attr("colorcode"));
	if (res != "task"){
		var initiativeId = obj.attr("InitiativeID")			
		var selected = Enumerable.From(c.AllInitiateSource()).FirstOrDefault(undefined, "$.InitiativeID == '" + initiativeId +"'")
		console.log(selected)
		// if (d.bdID != undefined && d.lcID != undefined){
		
		// selected.LifeCycleId = d.lcID;
		selected.LifeCycleId = target.attr("LCId");
		if (target.parent().attr("SCCategory") != 'All'){
			selected.SCCategory = target.parent().attr("SCCategory");
			selected.BusinessDriverId = target.parent().attr("BDId");
		}
		// }
		Initiative.Mode("copyclone");	
		Initiative.Add(selected);
	} else {
		var selected = Enumerable.From(c.DataSource().Data.TaskList).FirstOrDefault(undefined, "$.Id == '" + code +"'")
		// if (d.bdID != undefined && d.lcID != undefined){
		selected.LifeCycleId = target.attr("LCId");
		if (target.parent().attr("SCCategory") != 'All'){
			selected.SCCategory = target.parent().attr("SCCategory");
			selected.BusinessDriverId = target.parent().attr("BDId");
		}	
		// }
		Task.Add(selected);
	}

	$('#mdlcloneormove').modal('hide')
	
}

