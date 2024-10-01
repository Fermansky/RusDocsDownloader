<script setup>
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div>
    <img id="logo" alt="Wails logo" src="./assets/images/historyrussia-logo.png"/>
  </div>
  <div style="margin-top: 10px">
    俄罗斯历史档案获取工具
  </div>
  <HelloWorld/>
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
