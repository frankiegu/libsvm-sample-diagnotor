#diagnotor for libsvm formatted samples.

机器学习样本诊断工具，通过多个维度的数据统计，曝光样本的情况，以及判断样本是否“健康”。

##Usage

```
Usage of ./diagnotor:
        -cover-max="-1": Threshold for cover
        -cover-min="-1": Threshold for cover
        -enable-mi=false: Whether enable mutal infermation evaluation, this is expensive, default is false
        -feature-max="-1": Threshold for width
        -feature-min="-1": Threshold for width
        -group-tag="": feature group tags, seperated by comma
        -mutal-max="-1": Threshold for mutal
        -mutal-min="-1": Threshold for mutal
```

输入和输出的基本要求

* 程序读取stdin中的libsvm格式的样本行
* 程序默认使用threshold.json作为配置文件, 启动时会搜索同级目录下的该文件
* 程序的统计结果均写出到文本文件中
