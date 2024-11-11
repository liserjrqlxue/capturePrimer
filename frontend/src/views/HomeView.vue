<script setup lang="ts">
import { useI18n } from "vue-i18n";

import HelloWorld from "@/components/HelloWorld.vue";

import router from "../router";

const { t } = useI18n();

if (history.length <= 1) {
  router.replace("/capture");
}
</script>

<template>
  <div class="home">
    <!-- Logo -->
    <img class="logo" alt="Vue logo" src="../assets/logo.png" />
    <HelloWorld :msg="t('homepage.welcome')" />
    <!-- Bottom button -->
    <div class="link">
      <!-- <a
        href="https://wails.io/docs/gettingstarted/installation"
        class="btn start"
        >{{ t("homepage.getting-started") }}</a
      > -->
      <router-link class="btn start" to="/capture">{{
        t("homepage.capturePrimer")
      }}</router-link>
    </div>
  </div>
</template>

<style lang="scss">
.home {
  .logo {
    display: block;
    width: 500px;
    height: 500px;
    margin: 10px auto 10px;
  }
  .link {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
    justify-content: center;
    margin: 18px auto;
    .btn {
      display: block;
      width: 150px;
      height: 50px;
      line-height: 50px;
      padding: 0 5px;
      margin: 0 30px;
      border-radius: 8px;
      text-align: center;
      font-weight: 700;
      font-size: 16px;
      white-space: nowrap;
      text-decoration: none;
      cursor: pointer;
      &.start {
        background-color: #fd0404;
        color: #ffffff;
        &:hover {
          background-color: #ec2e2e;
        }
      }
      &.star {
        background-color: #ffffff;
        color: #fd0404;
        &:hover {
          background-color: #f3f3f3;
        }
      }
    }
  }
}
</style>
