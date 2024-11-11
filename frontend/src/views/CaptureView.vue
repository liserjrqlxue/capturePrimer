<script setup lang="ts">
import { useI18n } from "vue-i18n";

import { ref, computed } from "vue";
import { SelectFile } from "../../wailsjs/go/main/App";
import { RunCapture } from "../../wailsjs/go/main/App";
import { SaveRows } from "../../wailsjs/go/main/App";
import { ClipboardSetText } from "../../wailsjs/runtime/runtime";

const { t } = useI18n();

const filePath = ref("");
const primerLength = ref(40);
const looseness = ref(false);

const captureStatus = ref("");
const rows = ref<any[]>([]);

// Computed property to determine the class
const statusClass = computed(() => {
  return captureStatus.value === "fail" ? "text-red-500" : ""; // Tailwind class for red text
});
const copySuccess = ref(false);
const tooltipX = ref(0);
const tooltipY = ref(0);

const openFileDialog = async () => {
  try {
    const options = {
      Title: "选择文件",
      Filters: "",
    };
    filePath.value = await SelectFile(options.Title);
    console.log("Selected file path:", filePath.value);
    // 这里可以添加更多处理文件路径的逻辑
  } catch (error) {
    console.error("Error selecting file:", error);
  }
};

const runCapture = async () => {
  rows.value = [];
  await RunCapture(filePath.value, Number(primerLength.value), looseness.value)
    .then((allResult) => {
      console.log(allResult);
      captureStatus.value = allResult.status;
      if (allResult.results) {
        rows.value = convertResult2Table(allResult.results);
        console.log(rows);
      }
      document.getElementById("resultCard")?.classList.remove("hidden");
    })
    .catch((error) => {
      alert(error);
      console.log(error);
    });
};

const copyRows = (event: { clientX: number; clientY: number }) => {
  tooltipX.value = event.clientX;
  tooltipY.value = event.clientY - 30;
  ClipboardSetText(
    rows.value.map((row) => `${row.name}\t${row.seq}`).join("\n")
  )
    .then((success) => {
      if (success) {
        copySuccess.value = true;
        // 显示提示1秒后消失
        setTimeout(() => {
          copySuccess.value = false;
        }, 1000);
      }
    })
    .catch((error) => {
      alert(error);
      console.log(error);
    });
};

const saveRows = async () => {
  // const resultCard = document.getElementById('resultCard')
  // resultCard?.classList.add('hidden');
  console.log("save rows", rows.value);
  SaveRows(rows.value);
};

const convertResult2Table = (data: any) => {
  console.log(data);
  return data.flatMap(function (item: any) {
    const id = item.index;
    const st = item.status;
    const name = item.name;
    if (st == "success" && item.capturePrimers.length > 0) {
      return item.capturePrimers.flatMap(function (pair: any) {
        return [
          {
            index: id,
            status: st,
            name: pair.primer5F.name,
            seq: pair.primer5F.seq,
            start: pair.primer5F.start,
            end: pair.primer5F.end,
          },
          {
            index: id,
            status: st,
            name: pair.primer3R.name,
            seq: pair.primer3R.seq,
            start: pair.primer3R.start,
            end: pair.primer3R.end,
          },
        ];
      });
    } else {
      return [{ index: id, status: st, name: name, seq: st }];
    }
  });
};
</script>

<template>
  <div class="container mx-auto">
    <div class="mx-auto object-center py-4" style="--wails-draggable: no-drag">
      <div class="gap-4 p-4">
        <div class="flex w-full">
          <label class="w-1/6">{{ t("capturepage.selectFile") }}</label>
          <input
            class="w-5/6"
            type="text"
            @click="openFileDialog"
            v-model="filePath"
            required
          />
        </div>
        <div class="flex w-full items-center py-4">
          <label class="w-1/6">{{ t("capturepage.primerLength") }}</label>
          <input
            type="range"
            v-model="primerLength"
            min="17"
            max="60"
            class="w-4/6 slider h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
          />
          <label class="w-1/6 text-right mt-2">
            {{ t("capturepage.value") }}: {{ primerLength }}
          </label>
        </div>
        <div class="flex w-full">
          <label class="w-1/6">{{ t("capturepage.looseness") }}</label>
          <input type="checkbox" v-model="looseness" class="w-5/6" />
        </div>
        <div class="flex justify-end py-4">
          <button class="button w-5/6 bg-white" @click="runCapture">
            {{ t("capturepage.analysisBtn") }}
          </button>
        </div>
      </div>
    </div>
    <div
      id="resultCard"
      class="hidden gap-4 p-4"
      style="--wails-draggable: no-drag"
    >
      <div class="flex justify-between py-4">
        <div class="flex w-4/12 justify-between">
          <div>
            <span class="w-2/12"
              >{{ t("capturepage.resultCount") }}: {{ rows.length }}</span
            >
          </div>
          <div>
            <span class="w-2/12">
              {{ t("capturepage.resultStatus") }}:
              <span :class="statusClass">{{
                t("capturepage." + captureStatus)
              }}</span></span
            >
          </div>
        </div>
        <div class="flex w-1/6 justify-end">
          <button class="button px-2" @click="copyRows">
            {{ t("capturepage.copyBtn") }}
          </button>
          <!-- 成功提示，当 showSuccess 为 true 时显示 -->
          <div
            v-if="copySuccess"
            :style="{ left: tooltipX + 'px', top: tooltipY + 'px' }"
            class="absolute p-2 bg-gray-700 text-white rounded"
          >
            {{ t("capturepage.copySuccess") }}
          </div>
          <button class="button px-2 ml-2" @click="saveRows">
            {{ t("capturepage.saveBtn") }}
          </button>
        </div>
      </div>
      <table
        class="primerTable text-center border-collapse border-spacing-2 w-full border border-slate-400 dark:border-slate-500 dark:bg-slate-800"
      >
        <thead class="bg-slate-50 dark:bg-slate-700">
          <tr v-if="captureStatus === 'success'">
            <th>{{ t("capturepage.titleIndex") }}</th>
            <th>{{ t("capturepage.titleName") }}</th>
            <th>{{ t("capturepage.titleSequence") }}</th>
            <th>{{ t("capturepage.titleStart") }}</th>
            <th>{{ t("capturepage.titleEnd") }}</th>
          </tr>
          <tr v-else>
            <th>{{ t("capturepage.titleIndex") }}</th>
            <th>{{ t("capturepage.titleName") }}</th>
            <th>{{ t("capturepage.titleError") }}</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="row in rows">
            <tr v-if="row.status === 'success'">
              <td>{{ row.index }}</td>
              <td>{{ row.name }}</td>
              <td>{{ row.seq }}</td>
              <td>{{ row.start }}</td>
              <td>{{ row.end }}</td>
            </tr>
            <tr v-else class="text-red-500">
              <td>{{ row.index }}</td>
              <td>{{ row.name }}</td>
              <td>{{ row.status }}</td>
            </tr>
          </template>
        </tbody>
      </table>
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

.button {
  background-color: rgba(171, 126, 220, 0.9);

  &:hover {
    background-color: #d7a8d8;
    color: #ffffff;
  }
}

table {
  user-select: none;
  background-color: #d7a8d8;

  // th:nth-child(2),
  // th:nth-child(3),
  td:nth-child(2),
  td:nth-child(3) {
    user-select: text;
  }
}

thead tr {
  background-color: rgba(171, 126, 220, 0.9);

  &:hover {
    background-color: #d7a8d8;
    color: #ffffff;
  }
}

tbody {
  background-color: #d7a8d8;

  tr {
    background-color: #d7a8d8;

    &:hover {
      color: #ffffff;
    }
  }
}
</style>
