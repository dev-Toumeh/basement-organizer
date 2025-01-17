/** HandleUpdatePicture is used to update inputs for requests when a file is selected for updating a picture.
 * @param {HTMLInputElement} element
 * @param {Event} event
 * @param {String} targetElementID - target element should be HTMLInputElement. Will inform backend if picture should be updated with PUT request.
 *  */
function HandleUpdatePicture(element, event, targetElementID) {
    console.log(event);
    const btn = document.getElementById("remove-picture-button");
    if (element.files.length !== 0) {
        setInputChecked(targetElementID, true);
        setImgCrossedOut(false);
        btn.disabled = false;
    } else {
        setInputChecked(targetElementID, false);
    }
}

/** RemovePicture will inform backend to remove picture on PUT request.
 * @param {HTMLInputElement} srcElement - where the event comes from (should be a button).
 * @param {PointerEvent} event - should be click event
 * @param {String} pictureInputID - element where current picture information is stored, will be cleared.
 *  */
function RemovePicture(srcElement, event, pictureInputID) {
    console.debug("RemovePicture event: ", event);
    setInputChecked('update-picture', true)
    const pictureInput = document.getElementById(pictureInputID)
    pictureInput.value = '';
    setImgCrossedOut(true);
    createAndShowNotification("picture will be removed after update", NotificationTypeInfo)
    srcElement.disabled = true;
}

/**
 * @param {String} targetID - should be HTMLInputElement.
 * @param {Boolean} checked
 *  */
function setInputChecked(targetID, checked) {
    /** @type {HTMLInputElement | null} */
    const target = document.getElementById(targetID);
    if (target === null) {
        console.error("invalid input target ID " + targetElementID)
        return
    }
    if (checked === true) {
        target.setAttribute("checked", "checked");
        target.checked = true;
    } else {
        target.removeAttribute("checked");
        target.checked = false;
    }
}

/** 
 * @param {Boolean} crossedOut
 *  */
function setImgCrossedOut(crossedOut) {
    const img = document.getElementById("image-overlay");
    if (img === null) {
        return
    }
    if (crossedOut === true) {
        img.classList.add("image-crossed-out");
    } else {
        img.classList.remove("image-crossed-out");
    }
}
