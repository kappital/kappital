# Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM euleros:latest

WORKDIR /opt/kappital/manager

RUN groupadd -g 10000 kappital && \
    useradd -u 10000 -g 10000 kappital && \
    echo "Defaults targetpw" >> /etc/sudoers && \
    mkdir -p /opt/kappital/{manager,certs,log,audit,database} && \
    mkdir -p /opt/kappital/certs/conf

COPY bin/kappital-manager /opt/kappital/manager/

RUN chown -R kappital:kappital /opt/kappital/ && \
    chmod 750 /opt/kappital/manager && \
    chmod 550 /opt/kappital/manager/kappital-manager && \
    chmod -R 700 /opt/kappital/certs

