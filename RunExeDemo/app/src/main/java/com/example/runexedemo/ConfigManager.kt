package com.example.runexedemo

import android.content.Context
import java.io.File

object ConfigManager {
    private const val CONFIG_RELATIVE_PATH = "config/config.json"
    private const val DEFAULT_JSON = """
[
  {
    "QdInfo": "SO+aPyWTJ02k4C9FkkB29fACDXIsJx4pAGbhVI07D8hjHPOEsCFgpKs6sLj9MLaOMad0nLWkjLZiPGzrYzIs3aToDVNBVw9rru4owdN+fTUJ1mfNNfX+Cf/Jmpzp3FOyeRsYzY8OzU268PrSjfmoxkxOVShrJK1fFmvUCYngXO0/0f7+TVesNkPkn65tCM4Yre74sE5wA9/M3JNz8OTGmoLtMVjjG8AKtBGgwIHdRJ1ivLKft/7cX0YgTWLWHIK2rJJg0O+S5NZ6rEFO19M9QbzC5k5vWsDhBJOgRqFDOHgD9QNVopJwEQ==",
    "SdkSign": "fwU0VSlfsV9boLwdRxTi/h3OW8wGi5mhonPBvVq0SYpOFdUUSmDUkReMnIqD SCg0Vcnhn2BlfYbaLi+swH/kt21eudb8VJyBhfWdSJj9e1nhHCe2nEkV5evQ MwyNPLDDaBSTij/EvFalrvIY1oxJFkyLVGCMJncv2uqLPLw5mHmRtqp4G5Kv UeoK8Ylz7uD/xwyvrqpu4gY=",
    "YwKey": "ykAJPC3xEdkJ",
    "YwGuid": "120154865151",
    "Ibex": "vElkxzXlpCliiE5id5xnzGDwSGF05QBC-OWSXe_ZF3FXBoJj3cU25xfb4G4rkNslbz6Z9fF64KAeSvSdOZRJrDxDQWmtVl1hmZ0Fa-RagSK2t1lzWpUy41jmw_Yrc2oI7wodErXKAKB-KIhl16lbZ9IdAGeo_Tlr8d5cX_MzKkdYlygWpt3_I6_hExSwrqYz_GKBbAphtQDmDs-acU-pEmR1HYrBhTs2cIl5BqqKIhie9IFe-u6MAppgGMulwO0ba6ouem6xq12EXLQUrvyKup0Ll16j5rVljrHEwGbu5W1N2o7wtNS3Iydlw2vwuBowGCugIcv8qg8CfodN7mzDJvbraWbXsLq8dsSdDlyMVS0YkpME3G1AbznHCRKAkGk8ywbGaA9uOHtehmjzRGmW6DJZXjCEfaI4Sshz5TOZioqMMRV4SkLcXIeAMNIk2CIW03ANavCbkSo3MzQ1M2RiNjMxYjZmZTU0ZTU0ZjdjZDQyY2NmYzY5Zg==",
    "TaskType": [
      1,
      2,
      3,
      4
    ]
  },
  {
    "QdInfo": "SO+aPyWTJ02k4C9FkkB29fACDXIsJx4pAGbhVI07D8hjHPOEsCFgpKs6sLj9MLaOMad0nLWkjLZiPGzrYzIs3aToDVNBVw9rru4owdN+fTUJ1mfNNfX+CQBZaD0PyG/b1t2Pnsyhhg2tqkEvCBtHY4Fpbzx5ue7BlLnHwz2azU0q3PXPvqDJiiNs7vPWlE06k0BZFS9f5FLlGVUNhQIu1Nbd0D6O7DEuvvz+zZ2DSROns8VuJf8rzvxTR0ov/0tWLadDZDBYBwnybcKtnCv1s7H5yWbYARkcNm8B8FGmm24Bpa7hrHJcbg==",
    "SdkSign": "fwU0VSlfsV9boLwdRxTi/taXXOTQ560ia5UbjaWanjn4PYuj3MjbBFa158CR tem9o1nUAMjAJ+BZ9PRivJutqE22tXHM9j2ML/GKrUVYKs2rL64y/dWEI72d CQoZqG/fS49EqeSG5THmd/Sbb5YkHhxEo3hNmhGtqlE4Awemy9QOrVC+Od9x 5rICZ3onNAfHpKV/tR7sksM=",
    "YwKey": "ykPETBSHdTNS",
    "YwGuid": "460067960",
    "Ibex": "PPOUGCvmhVhyOSgyQwRF10bu1vfoZz1n6ToPvdqrqQC9AOI9BRkz--qG4H1wAKHv2BGftWpbiKkT83-nX18aUqDmdTfJxz6cGtievRl-WWWF3rgObhVb3txY9tPkX2785TiZ5qrBE8TV1sQdi8Kd2W0aNzM8nMgewIY3Mo3QYRexqlW-FM0pnBWpoTg3_tJIhZ5SWsQHxh3Thi__v-x3mgATMHuclSj_AHOc_U1SVBm_Bm6joqk8gVDXBGqCc57J2mCMQbV0tstDWdz7kbT2YZLwkZcWFf-eAv_wj6HVLV4ryVxc9mspBXLbpib-R3ZwDgOX-T70aq_qTpmc25uUACOGWL-6WHtoJ3CfWPPpG0HF3HhGolOdBvwNVhh0fCaN_cDqolDmJ2CbKpGWOYXYAonI02DzOqdvybaclX6Mgudb0aqhdnn5GDE1YzliZDNiYzZiNGE3OTgyMmUxMjBkZGY5ZTc0MTNh",
    "TaskType": [
      1,
      2,
      3
    ]
  }
]

    """

    fun getConfigFile(context: Context): File {
        return File(context.filesDir, CONFIG_RELATIVE_PATH)
    }

    fun ensureConfig(context: Context): File {
        val file = getConfigFile(context)
        if (!file.parentFile.exists()) file.parentFile.mkdirs()
        if (!file.exists()) {
            file.writeText(DEFAULT_JSON.trimIndent())
        }
        return file
    }

    fun readConfig(context: Context): String {
        val file = ensureConfig(context)
        return file.readText()
    }

    fun writeConfig(context: Context, content: String) {
        val file = getConfigFile(context)
        if (!file.parentFile.exists()) file.parentFile.mkdirs()
        file.writeText(content)
    }
}

