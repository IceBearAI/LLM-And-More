export function animate(
  el: Element,
  keyframes: Keyframe[] | PropertyIndexedKeyframes | null,
  options?: number | KeyframeAnimationOptions
) {
  if (typeof el.animate === "undefined") return { finished: Promise.resolve() };

  let animation: Animation;
  try {
    animation = el.animate(keyframes, options);
  } catch (err) {
    return { finished: Promise.resolve() };
  }

  if (typeof animation.finished === "undefined") {
    (animation as any).finished = new Promise(resolve => {
      animation.onfinish = () => {
        resolve(animation);
      };
    });
  }

  return animation;
}

export const doAnimation = {
  scaleIn(el: Element) {
    if (el) {
      animate(el, [{ transformOrigin: "center" }, { transform: "scale(1.03)" }, { transformOrigin: "center" }], {
        duration: 250,
        easing: "cubic-bezier(0.4, 0, 0.2, 1)"
      });
    }
  }
};
