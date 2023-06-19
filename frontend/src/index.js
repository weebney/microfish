import { GetPath } from "../wailsjs/go/main/App.js";
import { EventsEmit } from "../wailsjs/runtime/runtime.js";
import { EventsOn } from "../wailsjs/runtime/runtime.js";

let vf = document.getElementById("frame")

GetPath().then((path) => {
	path = encodeURIComponent(path)
	vf.src = "./src/pdfjs/web/viewer.html?file=" + path + "#pagemode=none"
});

EventsOn("reloadPage", () => {
	vf.contentWindow.location.reload();
})


vf.addEventListener('load', function () {
	var iframeDocument = vf.contentDocument || vf.contentWindow.document;
	
	// save/open workaround by cloning
	function replaceButtonWithClone(buttonId) {
		let oldBtn = iframeDocument.getElementById(buttonId);
		let newBtn = oldBtn.cloneNode(true);
		newBtn.id = oldBtn.id;  // ensure the cloned button has the same id
		oldBtn.parentNode.replaceChild(newBtn, oldBtn);
		return newBtn;
	}
	
	let newSaveBtn = replaceButtonWithClone("download");
	let newSaveBtn2 = replaceButtonWithClone("secondaryDownload");
	let newOpenBtn = replaceButtonWithClone("openFile");
	let newOpenBtn2 = replaceButtonWithClone("secondaryOpenFile");
	
	// add event listeners to the new buttons
	newSaveBtn.addEventListener("click", (e) => { 
		e.preventDefault();
		saveMe(); 
	});
	newSaveBtn2.addEventListener("click", (e) => {
		e.preventDefault();
		saveMe(); 
	});
	
	newOpenBtn.addEventListener("click", (e) => { 
		e.preventDefault();
		openMe(); 
	});
	newOpenBtn2.addEventListener("click", (e) => {
		e.preventDefault();
		openMe(); 
	});

	// this setTimeout waits til the pdf is loaded
	// what an ugly workaround... plsfix
	setTimeout(() => {
		var anchors = iframeDocument.getElementsByTagName('a');
		for (var i = 0; i < anchors.length; i++) {
			let a = anchors[i]
			let href = a.getAttribute('href');
			a.removeAttribute('href');

			// pretty sure this is broken
			if (href && href.indexOf('#') === 0) {
				continue;  // skip relative anchors
			}

			if (href && href.indexOf(iframeDocument.location.origin) !== 0) {
				a.addEventListener("click", () => { sendToBrowser(href) })
			}
		}

		/* // saving
		iframeDocument.addEventListener('keydown', e => {
			if ((e.ctrlKey || e.metaKey) && e.key === 's') {
				e.preventDefault();
				saveMe()
			}
		});

		let saveBtn = iframeDocument.getElementById("download")
		let saveBtn2 = iframeDocument.getElementById("secondaryDownload")

		saveBtn.addEventListener("click", () => { e.preventDefault(); saveMe() })
		saveBtn2.addEventListener("click", () => { e.preventDefault(); saveMe() })

		// opening 
		iframeDocument.addEventListener('keydown', e => {
			if ((e.ctrlKey || e.metaKey) && e.key === 'o') {
				e.preventDefault();
				openMe()
			}
		});

		let openBtn = iframeDocument.getElementById("openFile")
		let openBtn2 = iframeDocument.getElementById("secondaryOpenFile")

		openBtn.addEventListener("click", () => { e.preventDefault(); openMe() })
		openBtn2.addEventListener("click", () => { e.preventDefault(); openMe() }) */
	}, 2000);
});

function sendToBrowser(href) {
	EventsEmit("sendToBrowser", href)
}

function saveMe() {
	EventsEmit("saveDoc")
}

function openMe() {
	EventsEmit("openDoc")
}
