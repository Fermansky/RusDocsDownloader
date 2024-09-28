<script setup>
import {onMounted, ref} from 'vue'
import {GetPdf} from '../../wailsjs/go/main/App'
import {EventsOff, EventsOn, EventsOnce} from "../../wailsjs/runtime/runtime.js";

const startPage = ref(0)
const endPage = ref(0)
const fixedStartPage = ref(0)
const fixedEndPage = ref(0)
const bookName = ref("")
const quality = ref(8)
const processing = ref(false)

const progress = ref(0)

const getPagePreview = id => {
  if (!Number.isInteger(id)) {
    return ''
  }
  return `https://docs.historyrussia.org/pages/${id}/zooms/0`
}
const onPagePreviewLoaded = () => {

}

const sanitizeFileName = (name) => {
  return name.replace(/[<>:"\/\\|?*]+/g, '_')
}

function pageChanged() {
  const start = parseInt(startPage.value)
  const end = parseInt(endPage.value)

  fixedStartPage.value = start
  fixedEndPage.value = end

  if (end < start) {
    fixedEndPage.value = start
    endPage.value = fixedEndPage.value
  }
}

const download = () => {
  EventsOnce("complete", () => {
    processing.value = false
    EventsOff("progress")
  })
  EventsOn("progress", p => {
    progress.value = parseInt(p)
  })
  const sanitizedBookName = sanitizeFileName(bookName.value)
  GetPdf(parseInt(fixedStartPage.value), parseInt(fixedEndPage.value), sanitizedBookName)
  processing.value = true
  progress.value = 0
}

</script>

<template>
  <main>
    <div style="display: flex; justify-content: center">
      <div style="max-width: 120px">
        <div style="height: 150px">
          <el-image :src=getPagePreview(fixedStartPage) class="preview-page" @load="onPagePreviewLoaded">
            <template #error>
              <div/>
            </template>
          </el-image>
          <div>{{fixedStartPage}}</div>
        </div>
        <el-input :disabled="processing" v-model="startPage" type="number" min="0" step="1" @change="pageChanged" placeholder="请输入起始页码"/>
      </div>
      <div style="margin: 20px; min-width: 20px">
        <div v-if="fixedStartPage && fixedEndPage">{{fixedEndPage-fixedStartPage+1}}页</div>
      </div>
      <div style="max-width: 120px">
        <div style="height: 150px">
          <el-image :src=getPagePreview(fixedEndPage) class="preview-page" @load="onPagePreviewLoaded"/>
          <div>{{fixedEndPage}}</div>
        </div>
        <el-input :disabled="processing" v-model="endPage" type="number"
                  :min="fixedStartPage" step="1" @change="pageChanged" placeholder="请输入结束页码"/>
      </div>
    </div>
    <div style="margin: 20px">
      <el-input :disabled="processing" v-model="bookName" placeholder="请输入书名"></el-input>
    </div>
    <div style="padding: 0 20px">
      <div>
        图片质量
      </div>
      <el-slider :disabled="processing" v-model="quality" :step="1" :min="0" :max="8" show-stops show-input/>
    </div>
    <div>
      <el-button :disabled="processing" type="success" size="large" @click="download">下载</el-button>
    </div>
    <div v-if="processing" style="padding: 0 20px">
      <el-progress
          :percentage="progress"
          :stroke-width="15"
          striped
          striped-flow
          :duration="10"
      />
    </div>
  </main>
</template>

<style scoped>

.preview-page {
  height: 120px;
  min-height: 120px;
  max-width: 120px;
  min-width: 70px;
}

.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
