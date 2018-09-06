#! /bin/bash
# Verifies consistent generation, for client, server and models
# of a simple spec with various target packages params.
set -euxo pipefail
testcases="fixture-1683.yaml"

allTargets="target target-withDash target/subDir target-test target/sub1/sub2"
#for dir in ${allTargets}; do
#  rm -rf ${dir}*
#done
rm -rf codegen*

#targetDirs=${allTargets}
targetDirs="target"
#modelPkgs="mabc"
modelPkgs="mabc mabc/msubdir mabc-dashed mabc-dashed/msubdir mabc/sub1/msub2 mabc-test"
#serverPkgs="sabc"
serverPkgs="sabc sabc/ssubdir sabc-dashed sabc-dashed/ssubdir sabc/sub1/ssub2 sabc-test"
#clientPkgs="cabc"
clientPkgs="cabc cabc/csubdir cabc-dashed cabc-dashed/csubdir cabc/sub1/csub2 cabc-test"
#apiPkgs="aabc"
apiPkgs="aabc aabc/asubdir aabc-dashed aabc-dashed/asubdir aabc/sub1/asub2 aabc-test"
serverNames="nrcodegen"
#serverNames="nrcodegen nrcodegen-test nrcodegen_underscored" #no slashes to support here

opts=""
t=0
for spec in ${testcases} ; do
  log=${spec%.*}.log
  for target in ${targetDirs} ; do
    let t=t+1
    target="codegen${t}/${target}"
    echo "Testing target: ${target}"
    i=0
    for modelPkg in "" ${modelPkgs} ; do
      echo "Testing model package: ${modelPkg:-models}"
      let i=i+1
      tg=${target}${i}
      rm -rf ${tg} && mkdir -p ${tg}
      cmd="swagger generate model --skip-validation ${opts} --spec=${spec} --model-package=${modelPkg:-models} --target=${tg}"
      echo ${cmd} > ${tg}/README.md
      ${cmd} 1>${tg}/${log} 2>&1
      gen=${tg}/${modelPkg:-models}
      path=`dirname ${gen}`
      base=`basename ${gen}`
      base=`echo ${base}|tr '-' '_'`
      gen=${path}/${base}
      (cd ${gen}; go build && echo "Model OK")

      for serverPkg in "" ${serverPkgs} ; do
        echo "Testing server package: ${serverPkg:-restapi}"
        for serverName in ${serverNames} ; do
          for apiPkg in "" ${apiPkgs} ; do
            echo "Testing API package: ${apiPkg:-operations}"
            let i=i+1
            tg=${target}${i}
            rm -rf ${tg} && mkdir -p ${tg}
            cmd="swagger generate server --skip-validation ${opts} --spec=${spec} --model-package=${modelPkg:-models} --server-package=${serverPkg:-restapi} --api-package=${apiPkg:-operations} --target=${tg} --name=${serverName}"
            echo ${cmd} > ${tg}/README.md
            ${cmd} 1>${tg}/${log} 2>&1
            newName=$(echo ${serverName}|sed 's/-*test$//'|tr '_' '-')
            (cd ${tg}/cmd/${newName}"-server"; go build && echo "Server OK")
          done
        done
      done

      for clientPkg in "" ${clientPkgs} ; do
        echo "Testing client package: ${clientPkg:-client}"
        let i=i+1
        tg=${target}${i}
        rm -rf ${tg} && mkdir -p ${tg}
        cmd="swagger generate client --skip-validation ${opts} --spec=${spec} --model-package=${modelPkg:-models} --client-package=${clientPkg:-client} --target=${tg} --name=${serverName}"
        echo ${cmd} > ${tg}/README.md
        ${cmd} 1>${tg}/${log} 2>&1
        gen=${tg}/${clientPkg:-client}
        path=`dirname ${gen}`
        base=`basename ${gen}`
        base=`echo ${base}|tr '-' '_'`
        gen=${path}/${base}
        (cd ${gen} ; go build && echo "Client OK")
      done
    done
  done
done
