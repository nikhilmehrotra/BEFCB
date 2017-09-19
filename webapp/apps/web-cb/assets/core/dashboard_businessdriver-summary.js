var SummaryBusinessDriver = {
	Processing:ko.observable(false),
  Detail: ko.observable(),
  Total: ko.observable(0)
}

MetricsTypeList = ko.observableArray([
  {name: "Dollar Value ($)", value: "DOLLAR"},
  {name: "Numeric Value", value: "NUMERIC"},
  {name: "Percentage (%)", value: "PERCENTAGE"},
])
SummaryBusinessDriver.Remove = function(obj){
  SummaryBusinessDriver.Detail().BusinessMetrics.remove(obj);
}
function ActualDataAdd(ActualData, period, value, flag){
  var self = this;
  self = {
    Period:ko.observable(period),
    Value:ko.observable(value), 
    Flag:ko.observable(flag),
    Hapus:ko.observable("") 
  }

  self.Period.subscribe(function(newValue){
    current = _.find(ActualData, function (o) { return o.Flag == "C1"; });
    if(current == undefined){
      boolDateMustBeLowerThanCurrent = false;
    } else{
      if(current.Period == null || current.Period == ""){
        boolDateMustBeLowerThanCurrent = false;
      } else{
        boolDateMustBeLowerThanCurrent = newValue > current.Period;
      }
    }

    if(newValue == "" || newValue == null){
      self.Value(0)
      self.Hapus("hapus")
    }
    else if(boolDateMustBeLowerThanCurrent){
      self.Period(period)
      self.Value(value)
      swal('Warning', 'Date must be lower than Current Date', 'error')
    } else{
      tmpDate = _.find(ActualData, function(v){ return v.Period+"" === newValue+"" })
      if(tmpDate == undefined){
        //add to array
        self.Value(0)
      } else{
        self.Value(tmpDate.Value)
      }
    }
    // console.log(newValue, tmpDate)
  });

  return self;
}

SummaryBusinessDriver.Get = function(Id){
  data = _.find(c.DataSource().Data.SummaryBusinessDriver, function (o) { return o.Id == Id; });
  if(data != undefined){

    data = ko.mapping.toJS( data )

    _.each(data.BusinessMetrics, function(vv,ii){
      if(vv.ActualData != null){
        _.each(vv.ActualData, function(v,i){
          data.BusinessMetrics[ii].ActualData[i].Period = new Date(v.Period);
        })
      }

      data.BusinessMetrics[ii].TargetPeriod = (vv.TargetPeriod == "0001-01-01T00:00:00Z") ? "" : new Date(vv.TargetPeriod);
    })
    
    // SummaryBusinessDriver.Detail( data );
    SummaryBusinessDriver.Total( data.BusinessMetrics.length )
    _.each(data.BusinessMetrics, function(v,i){
      //ActualData
      Period = moment(new Date()).startOf('month')._d;
      Period1 = moment(Period). subtract(1, 'M').startOf('month')._d;
      Period2 = moment(Period1). subtract(1, 'M').startOf('month')._d;
      Period3 = moment(Period2). subtract(1, 'M').startOf('month')._d;

      if(v.ActualData == null){
        v.ActualData = [];

        v.ActualData.push(new ActualDataAdd(v.ActualData, Period, 0, 'C1'));
        v.ActualData.push(new ActualDataAdd(v.ActualData, Period1, 0, 'A3'));
        v.ActualData.push(new ActualDataAdd(v.ActualData, Period2, 0, 'A2'));
        v.ActualData.push(new ActualDataAdd(v.ActualData, Period3, 0, 'A1'));

        v.ActualData = _.sortBy(v.ActualData, 'Period')

      }
      else{

        // v.ActualData = _.sortBy(v.ActualData, 'Period')

        // //search month now
        // tmpDateNow = _.find(v.ActualData, function(v){ return v.Period+"" === Period+"" })
        // // console.log(tmpDateNow)
        // if(tmpDateNow === undefined){
        //   v.ActualData.push(new ActualDataAdd(v.ActualData, Period, 0, 'C1'));

        //   //search month before -1
        //   tmpDate1 = _.find(v.ActualData, function(v){ return v.Period+"" === Period1+"" })
        //   if(tmpDate1 === undefined){
        //     v.ActualData.push(new ActualDataAdd(v.ActualData, Period1, 0, 'A3'));

        //     //search month before -2
        //     tmpDate2 = _.find(v.ActualData, function(v){ return v.Period+"" === Period2+"" })
        //     if(tmpDate2 === undefined){
        //       v.ActualData.push(new ActualDataAdd(v.ActualData, Period2, 0, 'A2'));

        //       //search month before -3
        //       tmpDate3 = _.find(v.ActualData, function(v){ return v.Period+"" === Period3+"" })
        //       if(tmpDate3 === undefined){
        //         v.ActualData.push(new ActualDataAdd(v.ActualData, Period3, 0, 'A1'));
        //       }

        //     }

        //   }

        // }

        v.ActualData = _.sortBy(v.ActualData, 'Period')

        // console.log(v.ActualData)

        var a = [];
        var a1 = true
        var a2 = true
        var a3 = true
        var c1 = true
        _.each(v.ActualData, function(vv,ii){
          if(vv.Flag =="A1") a1 = false;
          if(vv.Flag =="A2") a2 = false;
          if(vv.Flag =="A3") a3 = false;
          if(vv.Flag =="C1") c1 = false;
          a.push(new ActualDataAdd(v.ActualData, vv.Period, vv.Value, vv.Flag))
        })
        if(a1) a.push(new ActualDataAdd(v.ActualData, null, 0, "A1"));
        if(a2) a.push(new ActualDataAdd(v.ActualData, null, 0, "A2"));
        if(a3) a.push(new ActualDataAdd(v.ActualData, null, 0, "A3"));
        if(c1) a.push(new ActualDataAdd(v.ActualData, null, 0, "C1"));
        v.ActualData = a;

      }

      //update All var
      data.BusinessMetrics[i] = v;
    })


    SummaryBusinessDriver.Detail(ko.mapping.fromJS( data ));

    if(data.BusinessMetrics == null || data.BusinessMetrics.length == 0){
      SummaryBusinessDriver.Add()
    }
    
    $("#summary-businessdriver").modal("show")
  }
}

SummaryBusinessDriver.Add = function(){
  Period = moment(new Date()).startOf('month')._d;
  Period1 = moment(Period). subtract(1, 'M').startOf('month')._d;
  Period2 = moment(Period1). subtract(1, 'M').startOf('month')._d;
  Period3 = moment(Period2). subtract(1, 'M').startOf('month')._d;

  var ActualData = [];
  ActualData.push({Period: Period, Value: 0, Flag: 'C1'});
  ActualData.push({Period: Period1, Value: 0, Flag: 'A3'});
  ActualData.push({Period: Period2, Value: 0, Flag: 'A2'});
  ActualData.push({Period: Period3, Value: 0, Flag: 'A1'});

  ActualData = _.sortBy(ActualData, 'Period')

  var a = [];
  _.each(ActualData, function(vv,ii){
    a.push(new ActualDataAdd(ActualData, vv.Period, vv.Value, vv.Flag))
  })
  ActualData = a;

  ActualData = ko.mapping.fromJS( ActualData )

  SummaryBusinessDriver.Detail().BusinessMetrics.push({
    DataPoint : ko.observable(""),
    MinLiabilities : ko.observable(0),
    MaxLiabilities : ko.observable(0),
    MinPrct : ko.observable(0),
    MaxPrct : ko.observable(0),
    MinValue : ko.observable(0),
    MaxValue : ko.observable(0),
    Completion : ko.observable(0),
    MetricType : ko.observable(""),
    ActualData : ActualData,
    TargetPeriod : ko.observable(Period),
    TargetValue : ko.observable(0)
  })
  SummaryBusinessDriver.Total(SummaryBusinessDriver.Total()+1)
}

SummaryBusinessDriver.Save = function(Id){
  var url = "/web-cb/dashboard/summarybusinessdriversave";

  var tmp = ko.mapping.toJS( SummaryBusinessDriver.Detail() )

  data = _.find(c.DataSource().Data.SummaryBusinessDriver, function (o) { return o.Id == tmp.Id; });
  data = ko.mapping.toJS( data );
  // if(data != undefined){
  //   data = ko.mapping.toJS( data );
  //   _.each(data.BusinessMetrics, function(v,i){
  //     tmpActualData = []
  //     _.each(v.ActualData, function(vv,ii){
  //       tmpActualData.push(vv)

  //       // a = _.find(tmp.BusinessMetrics[i].ActualData, function (o) { return o.Period+'' == new Date(vv.Period)+''; });
  //       // console.log('-->', a)

  //       a = true
  //       _.each(tmp.BusinessMetrics[i].ActualData, function(vvv,iii){
          
  //         if(new Date(vv.Period)+'' == vvv.Period+''){
  //           a = false
  //         }

  //       })

  //       if(a){
  //         tmpActualData.push()
  //         console.log('add new')
  //       }

  //     })

  //     console.log(tmpActualData)
  //   })
  // }

  if(data != undefined){

    _.each(tmp.BusinessMetrics, function(v,i){
      v.MinLiabilities = v.ActualData.length === 0 ? 0 : v.ActualData[0].Value;
      v.MinValue = v.ActualData.length === 0 ? 0 : v.ActualData[v.ActualData.length-1].Value;
      v.MaxLiabilities = v.TargetValue
      v.MaxValue = v.TargetValue
      var xx = v.MaxLiabilities + v.MinLiabilities;
      tmp.BusinessMetrics[i].MinPrct = v.MinValue/xx*100;
      tmp.BusinessMetrics[i].MaxPrct = 100;

      // a = [1,2,3]
      // b = [2,3,4]
      tmpActualData = []
      _.each(v.ActualData, function(vv,ii){
        if(vv.Period != null){
          tmpActualData.push(vv)
        }
      })

      if(data.BusinessMetrics[i] != undefined){
        _.each(data.BusinessMetrics[i].ActualData, function(vv,ii){
          a = _.find(v.ActualData, function (o) { return (o.Period+'' == new Date(vv.Period)+'' || o.Hapus == 'hapus'); });
          if(a == undefined){
            // console.log(vv)
            vv.Period = new Date(vv.Period)
            vv.Flag = ''
            tmpActualData.push(vv)
          }
        })
      }

      filterHapusArray = []
      _.each(tmpActualData, function(vv,ii){
        if(vv.Hapus != "hapus"){
          filterHapusArray.push(vv)
        }
      })

      tmp.BusinessMetrics[i].ActualData = filterHapusArray

      if(tmp.BusinessMetrics[i].TargetPeriod == null || tmp.BusinessMetrics[i].TargetPeriod == ""){
        tmp.BusinessMetrics[i].TargetPeriod = new Date("0001-01-01T00:00:00Z");
        tmp.BusinessMetrics[i].TargetValue = 0;
      }
    })  

    Param = {
      Id : tmp._id,
      Idx : tmp.Id,
      Name : tmp.Name,
      Seq : tmp.Seq,
      Type : tmp.Type,
      TargetPeriod: tmp.TargetPeriod,
      TargetValue: tmp.TargetValue,
      BusinessMetrics : tmp.BusinessMetrics,
      Parentid: tmp.Parentid,
      Parentname: tmp.Parentname,
      Category: tmp.Category
    };

    ajaxPost(url, Param, function (datas){
      if(datas == "success"){
        swal("Success", "Data Metric Success to save", "success");
      } else{
        swal("Warning", "Data Not Success to save", "error");
      }
      c.GetData();
    })
  } 

  $("#summary-businessdriver").modal("hide")
  $('.modal-backdrop').remove()
}