---
title: "Adopting Reason: strategies, dual sources of truth, and why genType is a big deal"
date: "2018-10-03"
tags:
  - "JavaScript"
  - "ReasonML"
  - "TypeScript"
  - "Flowtype"
  - "BuckleScript"
---

_Disclaimer: migrating an existing product to a new language, even if partially, is one of the most impactful decisions that a team or tech company might make, for better or worse. It is rarely justified, and in most cases the costs will outweigh the benefits. In any case, it is a decision that should be made after deep consideration and never without a large consensus or a detailed plan. Or in other words: do your own research._ 😉

---

![adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-01.jpeg](/media/adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-01.jpeg)

_"Alice knew which was which in a moment, because one of them had ‘DUM’ embroidered on his collar, and the other ‘DEE’”._

You have been looking at Reason and the value behind it:

- Rich type system with a [solid inference engine](https://codeburst.io/inference-engines-5-examples-with-typescript-flow-and-reason-edef2f4cf2d3).
- Very fast build times.
- Highly [performant](https://www.javierchavarri.com/performance-of-records-in-bucklescript/) and readable output code.
- Language-level [immutable](https://reasonml.github.io/docs/en/record) [data](https://reasonml.github.io/docs/en/list-and-array#list) [structures](https://reasonml.github.io/docs/en/tuple).

You start thinking Ken Wheeler is right: [Reason is serious business!](https://www.youtube.com/watch?v=lzEweA7RPi0)

After considering it, reading the documentation, discussing with colleagues, making some small apps from scratch, experiments and prototypes, you and your team / company decide its time to take a stab at migrating a small part of an existing app.

But… where does one start?

This article explores the two most common routes teams have followed in the past to progressively introduce a new language to an existing JavaScript codebase, and also explains how a recently released tool called [genType](https://github.com/cristianoc/genType) might become a key part of this process in the future.

Note: genType was originally called genFlow and the project [was renamed](https://github.com/cristianoc/genType/issues/40) recently.

## 🎩 Data-first strategy: the NoRedInk approach

In [a great article](https://blog.noredink.com/post/126978281075/walkthrough-introducing-elm-to-a-js-web-app), the team from NoRedInk detailed how they went about starting using Elm in their production app.

Their approach: start with the state management system (or business logic, [depending when you were born](https://twitter.com/dan_abramov/status/1025431204336742400) 😛).

This makes a lot of sense: considering how statically typed languages make runtime errors much harder to happen, migrating the parts of the code that deal with data seems like a great way to start getting a lot of value out of the initial investment. For NoRedInk, it was Flux stores and Elm. But the idea behind it can easily be ported to other state management system like Redux, or any other statically typed language, like Reason.

![adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-02.png](/media/adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-02.png)

*Migrating your app, one reducer at a time*

It allows to start simple. State management systems are quite contained, as they generally consist of pure functions that perform some updates to the incoming data. Leaving out side effects like DOM handling 🚁 or network management 📡 from the initial stages of the migration simplifies the process in many ways.

This approach has certainly benefits, but it comes with some constraints too:

- The new language data structures start being used, but this data is still needed in other parts of the app for some time until they are converted too (if that ever happens!). These parts expect the data to be in an idiomatic JavaScript format, so the new language must have good conversion and interop tools in order to keep the best developer experience.
- In the same vein, if the conversion between data types is slow or complex, this might lead to performance bottlenecks if the app has strong performance requirements.
- Lastly, but most importantly, data stores are generally used in many different parts of an app / product, so this approach probably requires buy-in from a larger group of stakeholders, and thus might be harder to adopt in larger teams or companies.

## 💬 UI-first strategy: the Messenger approach

For Facebook Messenger, the team decided to go with a different approach: start with the UI first, with the React components that had less data dependencies, and move “up” from there. This approach seems to have been successful, considering the reported news: [more than 50% of that codebase has been migrated to Reason](https://reasonml.github.io/blog/2017/09/08/messenger-50-reason).

The strategy has the great benefit of being much less intrusive than the data-first strategy described earlier. For larger teams –like the ones at Facebook– state management systems are generally quite ubiquitous, and can even be shared across apps. So being able to start adopting the language in isolation might be the only way forward initially, as it requires a much smaller amount of people to reach consensus: the team that is working with that part of the UI / app is enough.

![adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-03.png](/media/adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-03.png)

*It’s a long way to the top if you wanna rock’n’roll*

The downside is that it can take a while to use Reason for core data management pieces of the app, and this is one of the places where the language really excels. Because the migration is happening “from bottom to top”, the safety provided when dealing with core state, or handling any other critical data, might take longer to happen.

This approach also means you will have to deal with the DOM / UI from the get-go, which involves more complexity in the adoption process, a challenge represented by the “should I check the Reason or the ReasonReact docs?” question that Reason newcomers often struggle with, and that Keira Hodgkison covered in depth –together with other frictions– in her talk [“What’s not to love about Reason?”](https://www.youtube.com/watch?v=4xr0WE49eik) that I highly recommend.

To summarize, this UI-first approach is less demanding in human constraints (less stakeholders needed to take the decision) but adds more technical complexity, so the developers working on it will have to deal with different classes of problems (due to the DOM and UI libraries being part of the plan).

---

Both strategies are totally valid, and depending on your specific context you might want to go with one or another. But regardless which one you pick, there is a challenge that your team will face at some point.

## 🐘 The elephant in the room: dual sources of truth

The truth is, for any kind of migration of any app to a new language, there will be a lot of logic and data types that already exist in your app and are written in JavaScript, and you have to decide how to approach that.

There are many ways to tackle this in Reason, but at the higher level, we can separate them in two main directions:

#### 1. Consume existing JavaScript code from Reason, using bindings

There has been _a lot_ of work done to make writing bindings an easier process. BuckleScript has historically [provided](https://bucklescript.github.io/docs/en/interop-overview) many different ways to write bindings against almost any kind of JavaScript code. BuckleScript is always trying to remove divergences between data types, like [the recent changes to `option` types](https://bucklescript.github.io/blog/2018/07/17/release-4-0-0II) to make the compiled data more natural from the JavaScript side. There have also been additions like [`bs.abstract`](https://bucklescript.github.io/docs/en/object#record-mode), which allows to represent JavaScript objects almost seamlessly as Reason records.

But bindings are subject to break, as they introduce a secondary source of truth: **the bindings themselves hold one “truth”, but another one remains in the JavaScript code**.

Bindings for a public JavaScript library have less chances of breaking because they are public, and _in theory_ should be updated with more care. But in an internal app, contracts are more easily broken. And if a change in some JavaScript function is breaking the Reason code, it makes the investment in the migration harder to justify: we are introducing a statically typed language to make our app safer, but any boundary between the new and the existing language is unsafe.

![adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-04.png](/media/adopting-reason-strategies-dual-sources-of-truth-and-why-gentype-is-a-big-deal-04.png)

*Ce binding n’est pas what you think it is*

#### 2. Rewrite them in Reason, and expose them to JavaScript

The second option is to take the required functions and data types and progressively move them to Reason.

In case your app is written in TypeScript, or uses Flow, and you want to maintain the type safety in the consumers of these functions, you are also entering a dual source of truth problem: you would have to maintain Flow or TypeScript type definitions manually, and they would also be susceptible to break, if the Reason code changes at some point and you forget to update them.

---

This seems like a deal breaker. Does a migration always involve adding a secondary source of truth that has to be manually updated and thus introduces fragility? Or are there any alternatives?

Well, now there might be one!

## genType: Auto generation of idiomatic bindings between Reason and JavaScript

I recently started exploring [genType](https://github.com/cristianoc/genType), a tool created by [Jordan Walke](https://github.com/jordwalke) and [Cristiano Calcagno](https://github.com/cristianoc). genType reduces the impact of the dual source of truth problem by generating **high quality, idiomatic Flow or TypeScript type definitions from Reason code, ✨ automatically. ✨**

A whole lot of emphasis, as you can see, goes in

**✨ automatically. ✨**

One more time, for those in the back:

**✨✨✨ aaaaautomaticallyyyyyyyyy!! ✨✨✨**

Yes! It’s automatic, which is pretty awesome. And idiomatic. And high quality. All of it.

To be clear, **genType does not remove the dual source of truth**. As long as there are two languages involved, there will need to be two sources. But by automating the generation of boundaries between both sources, we remove the need of any human intervention to keep them in sync, which is what caused the trouble originally.

#### How does genType work?

To simplify a lot, genType allows you to add the annotation `[@genType]` to any declaration in Reason, like a record in this example:

```reason
[@genType]
type position = {
  latitude: int,
  longitude: int,
};
```

It will generate something like this in a co-located `.re.js` file (this is the Flow version, but it can produce a similar TypeScript output):

```flow
export type position = {|
  latitude: number,
  longitude: number
|};
```

Not only types, other values too! If a function uses the previously declared type, a converter to JS types will be generated too.

In `Coordinate.re` file, in Reason:

```reason
[@genType]
type position = {
  latitude: int,
  longitude: int,
};
[@genType]
let updateLatitude = p => {
  ...p,
  latitude: p.latitude + 2
};
```

Automatically generated `Coordinate.re.js` (a similar file can be obtained for TypeScript or vanilla JavaScript):

```flow
const CoordinateBS = require("./Coordinate.bs");
export type position = {|latitude:number, longitude:number|};
export const updateLatitude: position => position = function _(Arg1) {
  const result = CoordinateBS.updateLatitude([
    Arg1.latitude,
    Arg1.longitude
  ]);
  return {latitude: result[0], longitude: result[1]};
};
```

Notice how the converter:

- Picks the properties from the JavaScript object to call the generated function `CoordinateBS.updateLatitude` in `Coordinate.bs.js`.
- Then wraps them again in an object before returning them back to JavaScript.

This is a big deal! **It solves the problems originated from manually updating a secondary source of truth**, but also removes a bunch of other headaches that arise when making JavaScript and Reason work together in the same app:

- Makes sure any code migrated to Reason doesn’t break existing assumptions on the JavaScript side, as it complies with any existing Flow or TypeScript type definitions
- Allows to keep updating the Reason code, while providing a lot of guarantees that any relying JavaScript code is not breaking.
- Removes the need to pay a “developer experience tax” on either side. Now one can use native data structures and features in Reason (spread operator ftw!) and use objects or any other data type in JavaScript seamlessly, and with type safety.
- Applies to other areas besides state management, for example, it [facilitates the integration](https://www.youtube.com/watch?v=EV12EbxCPjM&feature=youtu.be) between JavaScript or TypeScript React components and ReasonReact components.
- Can potentially be used to expose your Reason library publicly so it can be consumed safely by JavaScript and TypeScript users with minimal costs of conversion.

#### What are the trade-offs?

The main trade-off of using genType is the potential impact in performance, as the converters come with their own runtime costs. If your app happens to require crossing the boundaries between Reason and JavaScript very often, then you might need to skip the converters generated by genType and fall back to the manual approach. However, given the compiler knowledge of all the shapes of records and objects, these conversions involve operations that are generally quite fast: reading object properties, making functions calls, or creating new objects and arrays. If your app uses structures with a large number of elements in them(object with many keys, arrays with many elements), it might need a separate solution that doesn’t involve converters.

---

genType is still in an early stage (v0.12.0 as of today, Oct 3), but it already solves a lot of different scenarios and it’s making fast progress to support more, with new releases happening almost daily. If you want to try it, you can follow the instructions in the [project repo](https://github.com/cristianoc/genType).

## Summing up

We have seen two strategies to migrate existing JavaScript apps to Reason progressively, and how genType can help doing so in a safely manner, without giving up ergonomics on either side.

If you have experiences migrating apps to Reason (or other languages!) or have suggestions or comments about the article, I look forward to your comments. Thanks for reading! ✌️

---

_Thanks to Cristiano Calcagno and Marcel Cutts for reviewing an early version of this article._
