<template>
  <div class="characters">
    <div
      ref="purpleRef"
      class="character purple"
      :style="purpleStyle"
    >
      <div class="eyes" :style="purpleEyesStyle">
        <div class="eye" :style="eyeBallStyle(18, isPurpleBlinking)">
          <div
            v-if="!isPurpleBlinking"
            class="pupil"
            :ref="setEyeRef('purple-left')"
            :style="pupilStyle('purple-left', 7, 5, purpleForce)"
          ></div>
        </div>
        <div class="eye" :style="eyeBallStyle(18, isPurpleBlinking)">
          <div
            v-if="!isPurpleBlinking"
            class="pupil"
            :ref="setEyeRef('purple-right')"
            :style="pupilStyle('purple-right', 7, 5, purpleForce)"
          ></div>
        </div>
      </div>
    </div>

    <div
      ref="blackRef"
      class="character black"
      :style="blackStyle"
    >
      <div class="eyes" :style="blackEyesStyle">
        <div class="eye" :style="eyeBallStyle(16, isBlackBlinking)">
          <div
            v-if="!isBlackBlinking"
            class="pupil"
            :ref="setEyeRef('black-left')"
            :style="pupilStyle('black-left', 6, 4, blackForce)"
          ></div>
        </div>
        <div class="eye" :style="eyeBallStyle(16, isBlackBlinking)">
          <div
            v-if="!isBlackBlinking"
            class="pupil"
            :ref="setEyeRef('black-right')"
            :style="pupilStyle('black-right', 6, 4, blackForce)"
          ></div>
        </div>
      </div>
    </div>

    <div
      ref="orangeRef"
      class="character orange"
      :style="orangeStyle"
    >
      <div class="eyes pupils-only" :style="orangeEyesStyle">
        <div
          class="pupil"
          :ref="setEyeRef('orange-left')"
          :style="pupilStyle('orange-left', 12, 5, orangeForce)"
        ></div>
        <div
          class="pupil"
          :ref="setEyeRef('orange-right')"
          :style="pupilStyle('orange-right', 12, 5, orangeForce)"
        ></div>
      </div>
    </div>

    <div
      ref="yellowRef"
      class="character yellow"
      :style="yellowStyle"
    >
      <div class="eyes pupils-only" :style="yellowEyesStyle">
        <div
          class="pupil"
          :ref="setEyeRef('yellow-left')"
          :style="pupilStyle('yellow-left', 12, 5, yellowForce)"
        ></div>
        <div
          class="pupil"
          :ref="setEyeRef('yellow-right')"
          :style="pupilStyle('yellow-right', 12, 5, yellowForce)"
        ></div>
      </div>
      <div class="mouth" :style="yellowMouthStyle"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'

type ForceLook = { x: number; y: number } | null

const props = defineProps<{ isTyping?: boolean; showPassword?: boolean; passwordLength?: number }>()

const purpleRef = ref<HTMLDivElement | null>(null)
const blackRef = ref<HTMLDivElement | null>(null)
const yellowRef = ref<HTMLDivElement | null>(null)
const orangeRef = ref<HTMLDivElement | null>(null)

const mouseX = ref(0)
const mouseY = ref(0)

const purplePos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })
const blackPos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })
const yellowPos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })
const orangePos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })

const isPurpleBlinking = ref(false)
const isBlackBlinking = ref(false)
const isLookingAtEachOther = ref(false)
const isPurplePeeking = ref(false)

const isHidingPassword = computed(() => (props.passwordLength || 0) > 0 && !props.showPassword)

let purpleBlinkTimer: ReturnType<typeof setTimeout> | null = null
let blackBlinkTimer: ReturnType<typeof setTimeout> | null = null
let peekTimer: ReturnType<typeof setTimeout> | null = null
let lookTimer: ReturnType<typeof setTimeout> | null = null

const clamp = (value: number, min: number, max: number) => Math.max(min, Math.min(max, value))

const updateCharacterPos = (el: HTMLDivElement | null, target: { faceX: number; faceY: number; bodySkew: number }) => {
  if (!el) return
  const rect = el.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const centerY = rect.top + rect.height / 3
  const deltaX = mouseX.value - centerX
  const deltaY = mouseY.value - centerY
  target.faceX = clamp(deltaX / 20, -15, 15)
  target.faceY = clamp(deltaY / 30, -10, 10)
  target.bodySkew = clamp(-deltaX / 120, -6, 6)
}

const handleMouseMove = (e: MouseEvent) => {
  mouseX.value = e.clientX
  mouseY.value = e.clientY
  updateCharacterPos(purpleRef.value, purplePos)
  updateCharacterPos(blackRef.value, blackPos)
  updateCharacterPos(yellowRef.value, yellowPos)
  updateCharacterPos(orangeRef.value, orangePos)
}

const scheduleBlink = (setter: (value: boolean) => void, assign: (timer: any) => void) => {
  const timeout = setTimeout(() => {
    setter(true)
    setTimeout(() => {
      setter(false)
      scheduleBlink(setter, assign)
    }, 150)
  }, Math.random() * 4000 + 3000)
  assign(timeout)
}

watch(
  () => props.isTyping,
  (val) => {
    if (val) {
      isLookingAtEachOther.value = true
      if (lookTimer) clearTimeout(lookTimer)
      lookTimer = setTimeout(() => {
        isLookingAtEachOther.value = false
      }, 800)
    } else {
      isLookingAtEachOther.value = false
    }
  }
)

watch(
  () => [props.passwordLength, props.showPassword],
  ([len, show]) => {
    if ((len || 0) > 0 && show) {
      const schedulePeek = () => {
        peekTimer = setTimeout(() => {
          isPurplePeeking.value = true
          setTimeout(() => {
            isPurplePeeking.value = false
          }, 800)
        }, Math.random() * 3000 + 2000)
      }
      schedulePeek()
    } else {
      if (peekTimer) clearTimeout(peekTimer)
      isPurplePeeking.value = false
    }
  }
)

onMounted(() => {
  window.addEventListener('mousemove', handleMouseMove)
  scheduleBlink((val) => (isPurpleBlinking.value = val), (t) => (purpleBlinkTimer = t))
  scheduleBlink((val) => (isBlackBlinking.value = val), (t) => (blackBlinkTimer = t))
})

onUnmounted(() => {
  window.removeEventListener('mousemove', handleMouseMove)
  if (purpleBlinkTimer) clearTimeout(purpleBlinkTimer)
  if (blackBlinkTimer) clearTimeout(blackBlinkTimer)
  if (peekTimer) clearTimeout(peekTimer)
  if (lookTimer) clearTimeout(lookTimer)
})

const purpleForce = computed<ForceLook>(() => {
  if ((props.passwordLength || 0) > 0 && props.showPassword) {
    return { x: isPurplePeeking.value ? 4 : -4, y: isPurplePeeking.value ? 5 : -4 }
  }
  if (isLookingAtEachOther.value) return { x: 3, y: 4 }
  return null
})

const blackForce = computed<ForceLook>(() => {
  if ((props.passwordLength || 0) > 0 && props.showPassword) return { x: -4, y: -4 }
  if (isLookingAtEachOther.value) return { x: 0, y: -4 }
  return null
})

const orangeForce = computed<ForceLook>(() => {
  if ((props.passwordLength || 0) > 0 && props.showPassword) return { x: -5, y: -4 }
  return null
})

const yellowForce = computed<ForceLook>(() => {
  if ((props.passwordLength || 0) > 0 && props.showPassword) return { x: -5, y: -4 }
  return null
})

const eyeRefs = new Map<string, HTMLElement>()
const setEyeRef = (key: string) => (el: HTMLElement | null) => {
  if (el) eyeRefs.set(key, el)
  else eyeRefs.delete(key)
}

const pupilStyle = (key: string, size: number, maxDistance: number, force: ForceLook) => {
  const el = eyeRefs.get(key)
  let x = 0
  let y = 0
  if (force) {
    x = force.x
    y = force.y
  } else if (el) {
    const rect = el.getBoundingClientRect()
    const centerX = rect.left + rect.width / 2
    const centerY = rect.top + rect.height / 2
    const deltaX = mouseX.value - centerX
    const deltaY = mouseY.value - centerY
    const dist = Math.min(Math.sqrt(deltaX ** 2 + deltaY ** 2), maxDistance)
    const angle = Math.atan2(deltaY, deltaX)
    x = Math.cos(angle) * dist
    y = Math.sin(angle) * dist
  }
  return {
    width: `${size}px`,
    height: `${size}px`,
    transform: `translate(${x}px, ${y}px)`
  }
}

const eyeBallStyle = (size: number, blinking: boolean) => ({
  width: `${size}px`,
  height: blinking ? '2px' : `${size}px`
})

const purpleStyle = computed(() => {
  const hiding = isHidingPassword.value
  const typing = props.isTyping
  const showPassword = props.showPassword
  let transform = `skewX(${purplePos.bodySkew}deg)`
  if (showPassword && (props.passwordLength || 0) > 0) {
    transform = 'skewX(0deg)'
  } else if (typing || hiding) {
    transform = `skewX(${purplePos.bodySkew - 12}deg) translateX(40px)`
  }
  return { transform, height: typing || hiding ? '440px' : '400px' }
})

const purpleEyesStyle = computed(() => {
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  const left = showPassword ? 20 : isLookingAtEachOther.value ? 55 : 45 + purplePos.faceX
  const top = showPassword ? 35 : isLookingAtEachOther.value ? 65 : 40 + purplePos.faceY
  return { left: `${left}px`, top: `${top}px` }
})

const blackStyle = computed(() => {
  const typing = props.isTyping
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  let transform = `skewX(${blackPos.bodySkew}deg)`
  if (showPassword) {
    transform = 'skewX(0deg)'
  } else if (isLookingAtEachOther.value) {
    transform = `skewX(${blackPos.bodySkew * 1.5 + 10}deg) translateX(20px)`
  } else if (typing || isHidingPassword.value) {
    transform = `skewX(${blackPos.bodySkew * 1.5}deg)`
  }
  return { transform }
})

const blackEyesStyle = computed(() => {
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  const left = showPassword ? 10 : isLookingAtEachOther.value ? 32 : 26 + blackPos.faceX
  const top = showPassword ? 28 : isLookingAtEachOther.value ? 12 : 32 + blackPos.faceY
  return { left: `${left}px`, top: `${top}px` }
})

const orangeStyle = computed(() => {
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  const transform = showPassword ? 'skewX(0deg)' : `skewX(${orangePos.bodySkew}deg)`
  return { transform }
})

const orangeEyesStyle = computed(() => {
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  const left = showPassword ? 50 : 82 + orangePos.faceX
  const top = showPassword ? 85 : 90 + orangePos.faceY
  return { left: `${left}px`, top: `${top}px` }
})

const yellowStyle = computed(() => {
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  const transform = showPassword ? 'skewX(0deg)' : `skewX(${yellowPos.bodySkew}deg)`
  return { transform }
})

const yellowEyesStyle = computed(() => {
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  const left = showPassword ? 20 : 52 + yellowPos.faceX
  const top = showPassword ? 35 : 40 + yellowPos.faceY
  return { left: `${left}px`, top: `${top}px` }
})

const yellowMouthStyle = computed(() => {
  const showPassword = props.showPassword && (props.passwordLength || 0) > 0
  const left = showPassword ? 10 : 40 + yellowPos.faceX
  const top = showPassword ? 88 : 88 + yellowPos.faceY
  return { left: `${left}px`, top: `${top}px` }
})
</script>

<style scoped lang="scss">
.characters {
  position: relative;
  width: 550px;
  height: 400px;
}

.character {
  position: absolute;
  bottom: 0;
  transition: all 0.7s ease-in-out;
}

.purple {
  left: 70px;
  width: 180px;
  background-color: #6c3ff5;
  border-radius: 10px 10px 0 0;
  z-index: 1;
  transform-origin: bottom center;
}

.black {
  left: 240px;
  width: 120px;
  height: 310px;
  background-color: #2d2d2d;
  border-radius: 8px 8px 0 0;
  z-index: 2;
  transform-origin: bottom center;
}

.orange {
  left: 0;
  width: 240px;
  height: 200px;
  background-color: #ff9b6b;
  border-radius: 120px 120px 0 0;
  z-index: 3;
  transform-origin: bottom center;
}

.yellow {
  left: 310px;
  width: 140px;
  height: 230px;
  background-color: #e8d754;
  border-radius: 70px 70px 0 0;
  z-index: 4;
  transform-origin: bottom center;
}

.eyes {
  position: absolute;
  display: flex;
  gap: 8px;
  transition: all 0.2s ease-out;
}

.eye {
  border-radius: 9999px;
  background-color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: all 0.15s ease;
}

.pupil {
  background-color: #2d2d2d;
  border-radius: 9999px;
  transition: transform 0.1s ease-out;
}

.pupils-only .pupil {
  background-color: #2d2d2d;
}

.mouth {
  position: absolute;
  width: 80px;
  height: 4px;
  background-color: #2d2d2d;
  border-radius: 9999px;
  transition: all 0.2s ease-out;
}
</style>
