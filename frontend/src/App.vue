<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { WindowMinimise, Quit } from "../wailsjs/runtime/runtime";

const { t, availableLocales: languages, locale } = useI18n();

locale.value = "zh-Hans";

const onclickLanguageHandle = (item: string) => {
  item !== locale.value ? (locale.value = item) : false;
};

const onclickMinimise = () => {
  WindowMinimise();
};

const onclickQuit = () => {
  Quit();
};

// document.body.addEventListener("click", function (event) {
//   // event.preventDefault();
// });
</script>

<template>
  <!-- Header -->
  <div class="header">
    <!-- navigation -->
    <div class="nav" style="--wails-draggable: drag">
      <router-link to="/">{{ t("nav.home") }}</router-link>
      <router-link to="/capture">{{ t("nav.capture") }}</router-link>
      <router-link to="/about">{{ t("nav.about") }}</router-link>
    </div>
    <!-- Menu -->
    <div class="menu">
      <div class="language">
        <div
          v-for="item in languages"
          :key="item"
          :class="{ active: item === locale }"
          @click="onclickLanguageHandle(item)"
          class="lang-item"
        >
          {{ t("languages." + item) }}
        </div>
      </div>
      <div class="bar">
        <div class="bar-btn" @click="onclickMinimise">
          {{ t("topbar.minimise") }}
        </div>
        <div class="bar-btn" @click="onclickQuit">{{ t("topbar.quit") }}</div>
      </div>
    </div>
  </div>
  <!-- Page -->
  <div class="view">
    <router-view />
  </div>
</template>

<style lang="scss">
@import url("./assets/css/reset.css");
@import url("./assets/css/font.css");

html {
  width: 100%;
  height: 100%;
}

body {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  font-family: "JetBrainsMono";
  background-color: transparent;
}

#app {
  position: relative;
  // width: 900px;
  // height: 520px;
  height: 100%;
  background-color: rgba(219, 188, 239, 0.9);
  overflow: hidden;
}

.header {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  align-items: center;
  justify-content: space-between;
  height: 50px;
  padding: 0 10px;
  background-color: rgba(171, 126, 220, 0.9);

  .nav {
    a {
      display: inline-block;
      min-width: 70px;
      height: 30px;
      line-height: 30px;
      padding: 0 5px;
      margin-right: 8px;
      background-color: #ab7edc;
      border-radius: 2px;
      text-align: center;
      text-decoration: none;
      color: #000000;
      font-size: 14px;
      white-space: nowrap;

      &:hover,
      &.router-link-exact-active {
        background-color: #d7a8d8;
        color: #ffffff;
      }
    }
  }

  .menu {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
    justify-content: space-between;

    .language {
      margin-right: 20px;
      border-radius: 2px;
      background-color: #c3c3c3;
      overflow: hidden;

      .lang-item {
        display: inline-block;
        min-width: 70px;
        height: 30px;
        line-height: 30px;
        padding: 0 5px;
        background-color: transparent;
        text-align: center;
        text-decoration: none;
        color: #000000;
        font-size: 14px;

        &:hover {
          background-color: #ff050542;
          cursor: pointer;
        }

        &.active {
          background-color: #ff050542;
          color: #ffffff;
          cursor: not-allowed;
        }
      }
    }

    .bar {
      display: flex;
      flex-direction: row;
      flex-wrap: nowrap;
      align-items: center;
      justify-content: flex-end;
      min-width: 150px;

      .bar-btn {
        display: inline-block;
        min-width: 80px;
        height: 30px;
        line-height: 30px;
        padding: 0 5px;
        margin-left: 8px;
        background-color: #ab7edc;
        border-radius: 2px;
        text-align: center;
        text-decoration: none;
        color: #000000;
        font-size: 14px;

        &:hover {
          background-color: #d7a8d8;
          color: #ffffff;
          cursor: pointer;
        }
      }
    }
  }
}

.view {
  position: absolute;
  top: 50px;
  left: 0;
  right: 0;
  bottom: 0;
  // overflow: hidden;
  overflow: auto;
}
</style>
