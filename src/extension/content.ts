export {}

console.log("Loaded content script")

// **** So we need to add all elements to the array, move the element

async function run() {
  window.scrollTo(0, document.body.scrollHeight)

  await new Promise((res) => setTimeout(res, 3000))

  const out = []

  for (const elem of Array.from(document.querySelector("[role=feed]").children)) {
    out.push(elem)
  }

  // **** We need some type of recursion where we have an event listener

  async function recurse(index: number) {
    if (index === out.length) return

    const elem = out[index]

    // Expand the text
    const msg = elem.querySelector("[data-ad-comet-preview=message]")

    if (msg) {
      const more = msg.querySelector("[role=button]")

      if (more) {
        // @ts-ignore
        more.click()
        await new Promise((res) => setTimeout(res, 100))
      }

      console.log(msg.textContent)
    }

    // Get the post URL
    const share = elem.querySelector('[aria-label="Send this to friends or post it on your Timeline."]')
    if (!!share) {
      // @ts-ignore
      share.click()
      await new Promise((res) => setTimeout(res, 1000))

      const shareOptionContainer = document.querySelector("[role=dialog]")
      if (!!shareOptionContainer) {
        const shareOptions = shareOptionContainer.querySelectorAll("[role=button]")

        if (!!shareOptions) {
          const copyLink = shareOptions[shareOptions.length - 1]

          copyLink.scrollIntoView()

          const url = await new Promise((res) => {
            copyLink.addEventListener("click", () => {
              navigator.clipboard.readText().then((url) => res(url))
            })
          })

          console.log(url)
          //   **** So now here we can take the values and submit them
          //   **** We need to store this in local storage to make sure no duplicates
          // **** We need some button where upon clicking it will send down
          // **** Also need a variable about the depth AND setting the API key
        }
      }
    }

    recurse(index + 1)
  }

  await recurse(0)
}

//     // @ts-ignore
//     // copyLink.click()
//     // await new Promise((res) => setTimeout(res, 1000))

//     // const url = await navigator.clipboard.readText()

//     // Get the text
//     if (msg && url) {
//       const more = msg.querySelector("[role=button]")

//       if (more) {
//         // @ts-ignore
//         more.click()
//         await new Promise((res) => setTimeout(res, 100))
//       }

//       console.log(msg.textContent)
//     }
// }

run()
