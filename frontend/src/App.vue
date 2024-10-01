<script setup>
import HelloWorld from './components/HelloWorld.vue'
import {CheckUpdate, Update} from "../wailsjs/go/main/App.js";
import {ref} from "vue";

const loading = ref(false)

const test = () => {
  loading.value = true
  Update()
}
</script>

<template>
  <div v-loading="loading">
    <div>
      <img id="logo" alt="Wails logo" src="./assets/images/historyrussia-logo.png"/>
    </div>
    <div style="margin-top: 10px" @dblclick="test">
      test-version
    </div>
    <HelloWorld/>
  </div>
</template>

<style>
#logo {
  display: block;
  width: 30%;
  height: 30%;
  margin: auto;
  padding: 5% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
