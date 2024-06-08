'use strict'

console.log('hello world')

const channels = {}
document.querySelectorAll('.pill').forEach(element => {
  const name = element.dataset.channelName
  channels[name] = true
})

function clickFilter(event) {
  const value = event.target.dataset.channelName
  console.log(`type ${event.target.dataset.filterType}`)
  console.log(`value ${value}`)

  channels[value] = !channels[value]

  updateChannels(channels)
}

document.querySelectorAll('.pill').forEach(element => {
  element.onclick = clickFilter
});

function updateChannels(channels) {
  document.querySelectorAll('.pill').forEach(filterPill => {
    const isOn = channels[filterPill.dataset.channelName]
    if (isOn) {
      filterPill.classList.remove('disabled')
    } else {
      filterPill.classList.add('disabled')
    }
  })

  document.querySelectorAll('.listing').forEach(listing => {
    const listingChannels = []
    listing.querySelectorAll('.pill').forEach(listingPill => {
      listingChannels.push(listingPill.dataset.channelName)
    })
    const anyEnabled = listingChannels.some(ch => channels[ch])
    listing.hidden = !anyEnabled
  })

  document.querySelectorAll('.day-section').forEach(day => {
    const dayListings = day.querySelectorAll('.listing')
    var hidden = true
    console.log(dayListings)
    dayListings.forEach(dayListing => {
      if (!dayListing.hidden) {
        hidden = false        
      }
    })
    day.hidden = hidden
  })
}
