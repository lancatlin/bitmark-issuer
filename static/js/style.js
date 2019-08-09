var tags=document.querySelectorAll('.feature')

console.log(tags)

function tagOn(event){
    let tag=event.currentTarget
    tag.style.animation="tagOn"
    tag.style.animationFillMode="forwards"
}

function tagOff(event){
    let tag=event.currentTarget
    tag.style.animation="tagOff 0.75s"
    tag.style.animationFillMode="forwards"
}

for(var item of tags){
    item.addEventListener('mouseover',tagOn)
    item.addEventListener('mouseout',tagOff)   
}